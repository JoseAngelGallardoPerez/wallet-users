print("Wallet Users")

load("ext://restart_process", "docker_build_with_restart")

cfg = read_yaml(
    "tilt.yaml",
    default = read_yaml("tilt.yaml.sample"),
)

local_resource(
    "users-build-binary",
    "make fast_build",
    deps = ["./cmd", "./internal", "./rpc/cmd", "./rpc/internal"],
)
local_resource(
    "users-generate-protpbuf",
    "make gen-protobuf",
    deps = ["./rpc/proto/users/users.proto"],
)

docker_build(
    "velmie/wallet-users-db-migration",
    ".",
    dockerfile = "Dockerfile.migrations",
    only = "migrations",
)
k8s_resource(
    "wallet-users-db-migration",
    trigger_mode = TRIGGER_MODE_MANUAL,
    resource_deps = ["wallet-users-db-init"],
)

wallet_users_options = dict(
    entrypoint = "/app/service_users",
    dockerfile = "Dockerfile.prebuild",
    port_forwards = [],
    helm_set = [],
)

if cfg["debug"]:
    wallet_users_options["entrypoint"] = "$GOPATH/bin/dlv --continue --listen :%s --accept-multiclient --api-version=2 --headless=true exec /app/service_users" % cfg["debug_port"]
    wallet_users_options["dockerfile"] = "Dockerfile.debug"
    wallet_users_options["port_forwards"] = cfg["debug_port"]
    wallet_users_options["helm_set"] = ["containerLivenessProbe.enabled=false", "containerPorts[0].containerPort=%s" % cfg["debug_port"]]

docker_build_with_restart(
    "velmie/wallet-users",
    ".",
    dockerfile = wallet_users_options["dockerfile"],
    entrypoint = wallet_users_options["entrypoint"],
    only = [
        "./build",
        "zoneinfo.zip",
    ],
    live_update = [
        sync("./build", "/app/"),
    ],
)
k8s_resource(
    "wallet-users",
    resource_deps = ["wallet-users-db-migration"],
    port_forwards = wallet_users_options["port_forwards"],
)

yaml = helm(
    "./helm/wallet-users",
    # The release name, equivalent to helm --name
    name = "wallet-users",
    # The values file to substitute into the chart.
    values = ["./helm/values-dev.yaml"],
    set = wallet_users_options["helm_set"],
)

k8s_yaml(yaml)

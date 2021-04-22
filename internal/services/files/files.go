package files

import (
	"github.com/Confialink/wallet-users/internal/srvdiscovery"
	"context"
	"errors"
	"net/http"

	"github.com/inconshreveable/log15"

	pb "github.com/Confialink/wallet-files/rpc/files"
)

type FilesService struct {
	filesProcessor pb.ServiceFiles
	logger         log15.Logger
}

func NewFilesService(logger log15.Logger) *FilesService {
	return &FilesService{logger: logger}
}

func (s *FilesService) UserHasFiles(uid string, excludeCategories []string) (resp *pb.UserHasFilesResp, err error) {
	if s.processor() == nil {
		return resp, errors.New("can't connect to files")
	}

	request := pb.UserHasFilesReq{Uid: uid, ExcludeCategories: excludeCategories}
	resp, err = s.processor().UserHasFiles(context.Background(), &request)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FilesService) Upload(
	bytes []byte,
	fileName string,
	uid string,
	adminOnly bool,
	private bool,
	category string,
) (resp *pb.UploadFileResp, err error) {
	if s.processor() == nil {
		return resp, errors.New("can't connect to files")
	}

	request := pb.UploadFileReq{
		Bytes:     bytes,
		Uid:       uid,
		FileName:  fileName,
		AdminOnly: adminOnly,
		Private:   private,
		Category:  category,
	}
	resp, err = s.processor().UploadFile(context.Background(), &request)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FilesService) processor() pb.ServiceFiles {
	if s.filesProcessor == nil {
		connection, err := getFilesClient()
		if err != nil {
			s.logger.Error("Failed to connect to files")
			return nil
		}

		s.filesProcessor = connection
	}

	return s.filesProcessor
}

func getFilesClient() (pb.ServiceFiles, error) {
	filesUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameFiles)
	if err != nil {
		return nil, err
	}
	return pb.NewServiceFilesProtobufClient(filesUrl.String(), http.DefaultClient), nil
}

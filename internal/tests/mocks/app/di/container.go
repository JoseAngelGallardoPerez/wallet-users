package di

import (
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/jinzhu/gorm"
)

func MockConnectionInContainer(gormDb *gorm.DB) {
	//gormDb.LogMode(true)
	/* container := di.Container
	pointerVal := reflect.ValueOf(container)
	val := reflect.Indirect(pointerVal)
	connectionFiled := val.FieldByName("dbConnection")
	ptrToConnectionFiled := unsafe.Pointer(connectionFiled.UnsafeAddr())
	realPtrToConnectionFiled := (**gorm.DB)(ptrToConnectionFiled)
	*realPtrToConnectionFiled = gormDb */
}

func MockSysSettingsInContainer(settings *syssettings.SysSettings) {
	/* container := di.Container
	pointerVal := reflect.ValueOf(container)
	val := reflect.Indirect(pointerVal)
	connectionFiled := val.FieldByName("sysSettings")
	ptrToConnectionFiled := unsafe.Pointer(connectionFiled.UnsafeAddr())
	realPtrToConnectionFiled := (**syssettings.SysSettings)(ptrToConnectionFiled)
	*realPtrToConnectionFiled = settings */
}

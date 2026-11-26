package manager

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"github.com/ricochet2200/go-disk-usage/du"
)

var KB = uint64(1024)

type Service interface {
	CheckHealth() ManagerHealth
}

func NewManagementService() Service {
	return ManagementService{}
}

type ManagementService struct {
}

func (m ManagementService) CheckHealth() ManagerHealth {

	hDB, upDB := m.healthDB()
	hDisk, upDisk := m.healthDiskUsage()
	status := "UP"

	if !upDB || !upDisk {
		status = "DOWN"
	}

	return ManagerHealth{
		Status:    status,
		DB:        hDB,
		DiskSpace: hDisk,
	}

}

func (m ManagementService) healthDB() (ManagerHealthDB, bool) {

	hello, errDB := tenants.TenantService.PingMainTenant()
	statusDB := "UP"
	if errDB != nil {
		statusDB = "DOWN"
	}

	dbHealth := ManagerHealthDB{
		Status: statusDB,
		Hello:  hello.Result,
	}

	go tenants.TenantService.LogOpenConnections()

	return dbHealth, statusDB == "UP"
}

func (m ManagementService) healthDiskUsage() (ManagerHealthDiskSpace, bool) {

	usage := du.NewDiskUsage("/")
	statusDisk := "UP"

	if (usage.Usage() * 100) > 90 {
		statusDisk = "DOWN"
	}

	diskSpace := ManagerHealthDiskSpace{
		Free:   usage.Free() / (KB * KB),
		Total:  usage.Size() / (KB * KB),
		Status: statusDisk,
	}

	return diskSpace, statusDisk == "UP"
}

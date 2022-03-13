package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(connectionString string) (*PostgresRepo, error) {
	var err error
	pr := new(PostgresRepo)
	pr.db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, err
	}
	pr.db.AutoMigrate(&Download{}, &Hub{}, &DownloadHub{})
	return pr, nil
}

func (p *PostgresRepo) Insert(model interface{}) error {
	result := p.db.Create(model)
	return result.Error
}

func (p *PostgresRepo) Update(model interface{}, data map[string]interface{}) error {
	result := p.db.Model(model).Updates(data)
	return result.Error
}

func (p *PostgresRepo) GetFirst(model interface{}, id uint) error {
	res := p.db.First(model, id)
	return res.Error
}

func (p *PostgresRepo) Get(model interface{}, data map[string]interface{}) error {
	res := p.db.Model(model).Where(data).Find(model)
	return res.Error
}

func (p *PostgresRepo) GetDownloadFromHub(model *Download, hubId uint, downloadId uint) error {
	res := p.db.Table("download").Select(`download.eta, 
										 download.download_id,
										 download.speed, 
										 download.total, 
										 download.progress, 
										 download.url, 
										 download.type, 
										 download.format, 
										 download.dir, 
										 download.details`).Joins("INNER JOIN download_hub ON download.download_id = download_hub.download_id").
		Where("download_hub.hub_id = ? and download_hub.download_id= ?", hubId, downloadId).Scan(&model)
	return res.Error
}

func (p *PostgresRepo) GetDownloadsFromHub(model *[]Download, hubId uint) error {
	res := p.db.Table("download").Select(`download.eta, 
										 download.download_id,
										 download.speed, 
										 download.total, 
										 download.progress, 
										 download.url, 
										 download.type, 
										 download.format, 
										 download.dir, 
										 download.details`).Joins("INNER JOIN download_hub ON download.download_id = download_hub.download_id").
		Where("download_hub.hub_id = ?", hubId).Scan(&model)
	return res.Error
}

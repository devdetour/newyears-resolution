package contracts

import "github.com/devdetour/ulysses/models"

type IContractStore interface {
	StoreContract(models.BasicContract)
	GetContractById(uint) models.BasicContract
	GetAllContracts() []models.BasicContract
}

type SimpleContractStore struct {
	contracts map[uint]models.BasicContract
}

func (s *SimpleContractStore) GetAllContracts() {
	v := make([]models.BasicContract, 0, len(s.contracts))

	// todo - maybe hold on to list we currently have if one has not been added.
	for _, value := range s.contracts {
		v = append(v, value)
	}
}

func (s *SimpleContractStore) GetContractById(id uint) models.BasicContract {
	return s.contracts[id]
}

func (s *SimpleContractStore) StoreContract(c models.BasicContract) {
	s.contracts[c.UserID] = c
}

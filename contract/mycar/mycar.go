package main

//1. import
import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"log"


	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
//2. chaincode 구조체
 type SmartContract struct {
	contractapi.Contract
 }
//3. car, query result 구조체
type Car struct {
	Make string `json:"make"`
	Model string `json:"model"`
	Colour string `json:"colour"`
	Owner string `json:"owner`
}
type QueryResult struct {
		Key string `json:"Key"`
		Record *Car
	}

//4.1 initLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour:"blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour:"red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour:"green", Owner:"Jin Soo"},
		Car{Make:"Volkswagen", Model: "Passat", Colour:"yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner:"Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: " indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}


	for	i, car := range cars {
			carAsBytes, _ := json.Marshal(car)
			err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

			if err != nil {
				return fmt.Errorf("Failed to put to world state. %s", err.Error())
			}
		}

		return nil
	}

	type HistoryQueryResult struct {
		Record    *Car      `json:"record"`
		TxId      string    `json:"txId"`
		Timestamp time.Time `json:"timestamp"`
		IsDelete bool       `json:"isDelete"`
	}

	func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, carID string )([]HistoryQueryResult, error) {
		log.Printf("GetAssetHistory: Id %v", carID)

		resultsIterator, err := ctx.GetStub().GetHistoryForKey(carID)
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()

		var records []HistoryQueryResult
		for resultsIterator.HasNext() {
			response, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			} 

			var car Car
			if len(response.Value) > 0 {
				err= json.Unmarshal(response.Value, &car)
				if err != nil {
					return nil, err
				}
			}	else {
					car =Car{}
				}
			

			timestamp, err := ptypes.Timestamp(response.Timestamp)
			if err != nil {
				return nil, err
			}

			record := HistoryQueryResult{
				TxId:      response.TxId,
				Timestamp: timestamp,
				Record:    &car,
				IsDelete:  response.IsDelete,
			}
			records = append(records, record)
		}
		return records, nil
	}
//4.2 CreateOrder
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make: make,
		Model: model,
		Colour: colour,
		Owner: owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}
//4.3 QueryCar
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	
	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}
//4.4 ChangeCarOwner
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}
	
	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}
//4.5 QueryAllCars
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key : queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil


}
//5. main
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
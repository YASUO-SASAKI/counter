package main

import(
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CounterChaincode struct {
}

// カウンター情報
type Counter struct {
	Name string `json:"name"`
	Counts int32 `json:"counts"`
}

const numOfCounters int=3

// カウンター情報の初期値を設定
func (cc *CounterChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var counters [numOfCounters]Counter
	var countersBytes [numOfCounters][]byte

	// カウンター情報を生成
	counters[0] = Counter{Name: "Office Worker", Counts:0}
	counters[1] = Counter{Name: "Home Worker", Counts:0}
	counters[2] = Counter{Name: "Student", Counts:0}

	//カウンター情報をわーるどすてーとについか
	for i := 0; i < len(counters); i++ {
		//JSON形式に変換
		countersBytes[i], _ = json.Marshal(counters[i])
		//ワールドステートに追加
		stub.PutState(strconv.Itoa(i), countersBytes[i])
	}

	return nil, nil
}

// カウンター情報を更新
func (cc *CounterChaincode) Invoke(stub *shim.ChaincodeStub,function string,args []string) ([]byte, error) {
	// function名でハンドリング
	if function == "countUp" {
		// カウントアップを実行
		return cc.countUp(stub,args)
	}

	return nil, errors.New("Received unknown function")
}

// カウンター情報を参照
func (cc *CounterChaincode) Query(stub *shim.ChaincodeStub,function string,args []string) ([]byte,error) {
	// function名でハンドリング
	if function == "refresh" {
		// カウンター情報を崇徳
		return cc.getCounters(stub,args)
	}

	return nil, errors.New("Received unknown function")
}

// カウントアップを実行
func (cc *CounterChaincode) countUp(stub *shim.ChaincodeStub, args []string) ([]byte, error){
	// ワールドステートから選択されたカウンター情報を取得
	counterId := args[0]
	counterJson, _ := stub.GetState(counterId)

	// 取得したJSON形式の情報をCounterに変換
	counter := Counter{}
	json.Unmarshal(counterJson, &counter)

	// カウントアップ
	counter.Counts++

	// ワールドステートに更新後の値を追加
	counterJson, _ = json.Marshal(counter)
	stub.PutState(counterId,counterJson)

	return nil,nil
}

// カウンター情報の取得
func (cc *CounterChaincode) getCounters(stub *shim.ChaincodeStub,args []string) ([]byte, error) {
	var counters [numOfCounters]Counter
	var countersBytes [numOfCounters][]byte

	for i := 0; i < len(counters); i++ {
		// カウンター情報をワールドステートから取得
		countersBytes[i], _ = stub.GetState(strconv.Itoa(i))

		// 取得したJSON形式の情報をCounterに変換
		counters[i] = Counter{}
		json.Unmarshal(countersBytes[i], &counters[i])
	}

	//json形式に変換
	return json.Marshal(counters)
}

// Validating Peerに接続し、チェーンコードを実行suru
func main() {
	err := shim.Start(new(CounterChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s",err)
	}
}
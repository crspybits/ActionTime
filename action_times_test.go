package actiontimes

import "testing"
import "encoding/json"
import "time"
import "fmt"

func addAction(t *testing.T, s string) {
	err := AddAction(s)
	if err != nil {
		errorString := fmt.Sprintf("%#v", err)
		t.Error("Error adding action: " + s + "; " + errorString)
	}
}

func TestAddActionFailures(t *testing.T) {
	testData := []string{
		`foo`,
		`{}`,
		`{"foo": 123}`,
		`{"action": "jump"}`,                 // missing time field
		`{"time": 100}`,                      // missing action field
		`{"time": "jump", "action": "jump"}`, // bad type for time field
		`{"time": 100, "action": 100}`,       // bad type for action field
	}

	for _, data := range testData {
		err := AddAction(data)
		if err == nil {
			t.Error("Expected json parsing failure: " + data)
		}
	}
}

func getStats(t *testing.T) string {
	s := GetStats()
	if s == "" {
		t.Error("Error calling GetStats")
	}
	return s
}

func checkAction(t *testing.T, results []OutputAction, name string, expectedValue float32) {
	for _, result := range results {
		if result.Action == name {
			if result.Avg != expectedValue {
				t.Error("Action with name: " + name + " did not have expected value: " + fmt.Sprintf("%f", expectedValue) +
					"; actual value was: " + fmt.Sprintf("%f", result.Avg))
			}
			return
		}
	}

	t.Error("No action with name: " + name)
}

func TestSpecificKeyNamesAddActionSuccess(t *testing.T) {
	Reset()

	addAction(t, `{"action":"run", "time":100}`)
	s := getStats(t)
	expectedResult := `[{"action":"run","avg":100}]`
	if s != expectedResult {
		t.Error("Did not get expected key/values: " + s)
	}
}

func TestMultipleDifferentAddActionSuccess(t *testing.T) {
	Reset()

	addAction(t, `{"action":"jump", "time":100}`)
	addAction(t, `{"action":"jump", "time":200}`)
	addAction(t, `{"action":"run", "time":75}`)
	addAction(t, `{"action":"bling", "time":800}`)

	s := getStats(t)
	// fmt.Print("json: " + s + "\n")

	var data []OutputAction

	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		t.Error("Failure parsing json string results: " + s)
	}

	checkAction(t, data, "jump", 150)
	checkAction(t, data, "run", 75)
	checkAction(t, data, "bling", 800)
}

func TestMultipleSameAddActionSuccess(t *testing.T) {
	Reset()

	times := []float32{100, 130, 20, 50, 80, 245}

	total := float32(0.0)
	for _, v := range times {
		total += v
	}

	average := total / float32(len(times))
	action := "jump"

	for _, v := range times {
		addAction(t, "{\"action\":\""+action+"\", \"time\":"+fmt.Sprintf("%f", v)+"}")
	}

	s := getStats(t)
	// fmt.Print("json: " + s + "\n")

	var data []OutputAction

	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		t.Error("Failure parsing json string results: " + s)
	}

	checkAction(t, data, "jump", average)
}

func workerMethod(t *testing.T, action string, expectedValue float32) {
	addAction(t, "{\"action\":\""+action+"\", \"time\":"+fmt.Sprintf("%f", expectedValue)+"}")

	s := getStats(t)
	// fmt.Print("json: " + s + "\n")

	var data []OutputAction

	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		t.Error("Failure parsing json string results: " + s)
	}

	checkAction(t, data, action, expectedValue)
}

func worker(t *testing.T, action string, expectedValue float32) {
	for i := 0; i < 300; i++ {
		workerMethod(t, action, expectedValue)
	}
}

// This test is a little weak. I don't have a really good way of testing for
// renentrancy of the AddAction/GetStats methods.
func TestParallelAddActionSuccesses(t *testing.T) {
	Reset()

	go worker(t, "jump", 100)
	go worker(t, "run", 50)

	// So workers don't terminate early.
	time.Sleep(time.Second * 60)
}

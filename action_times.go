package actiontimes

import "encoding/json"
import "errors"
import "sync"
// import "fmt" // for debugging

// An action name and time parsed from JSON in input data.
// Using pointers in order to detect missing data
type inputaction struct {
  // names of fields must be those as given as they reflect keys in JSON input.
  // The field names must be uppercased so deserialization will find them in input.
  Action *string // action name
  Time *float32 // action duration
}

// An action name and average duration used in JSON in output data.
type OutputAction struct {
  // names of fields must be those as given as they reflect keys in JSON output.
  Action string // action name
  Avg float32 // average action duration
}

// I'm going to track an incremental running average for each action.
// In this way, I won't have to keep track of all times for each action which could
// amount to a lot of data over time.
type current_average struct {
  average float32
  n int
}

// Add next value into the incremental average
func (r *current_average) addValue(value float32) {
  // https://math.stackexchange.com/questions/106700/incremental-averageing
  r.n += 1
  r.average = r.average + (value - r.average)/float32(r.n)

  // fmt.Printf("addValue: average[%f]\n", r.average)
}

// I'm going to keep the data for the incremental averages in a map from action name
// to current average. Note that this is a static (not instance) variable. The requirement is for
// a library to keep track of averages, which suggests just a single set of averages.
var averages = make(map[string]*current_average)

// To establish a critical section around averages
var mutex = &sync.Mutex{}

// Accepts a JSON string of the following form:
//    {"action":"jump", "time":100}
// and maintains an average time for each action.
// If the input JSON cannot be parsed, or doesn't contain the expected key/values,
// this returns an error. No error is returrned if extra key/values are in the string.
func AddAction(s string) error {
  var time inputaction
  err := json.Unmarshal([]byte(s), &time)
  if err != nil {
    return err
  }

  if time.Action == nil || time.Time == nil {
    return errors.New("Missing required field in JSON")
  }

  mutex.Lock()
  _, exists := averages[*time.Action]
  if !exists {
    averages[*time.Action] = new(current_average)
  }

  averages[*time.Action].addValue(*time.Time)
  mutex.Unlock()

  // fmt.Printf("AddAction: key[%s] value[%f]\n", *time.Action, *time.Time)

  return nil // no error
}

// Returns an empty string if there was an error.
func GetStats() string {
  var results = []OutputAction{}

  // fmt.Printf("Number of actions: %d\n", len(averages))

  mutex.Lock()
  for action, average := range averages {
    // fmt.Printf("key[%s] value[%f]\n", action, average.average)

    // Passing &(action) directly into the ActionTime struct results in a bug. I think it's a string
    // allocation issue. Instead, do it this way:
    // var name = action

    next := OutputAction{action, average.average}
    results = append(results, next)
  }
  mutex.Unlock()

  bytes, err := json.Marshal(results)
  if err != nil {
    return ""
  }

  return string(bytes)
}

// This was not in the requirements but was needed for testing. It may also
// be useful in a practical application.
func Reset() {
  mutex.Lock()
  averages = make(map[string]*current_average)
  mutex.Unlock()
}

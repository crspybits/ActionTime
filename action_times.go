package actiontimes

import "encoding/json"
import "errors"
import "sync"
//import "fmt" // for debugging

// Using pointers in order to detect missing data
// Using this for both input (time is a specific time), and output (time is an average).
type ActionTime struct {
  Action *string
  Time *float32
}

// I'm going to track an incremental running average for each action.
// In this way, I won't have to keep track of all times for each action which could
// amount to a lot of data over time.
type currentAverage struct {
  Average float32
  N int
}

// Add next value into the incremental average
func (r *currentAverage) addValue(value float32) {
  // https://math.stackexchange.com/questions/106700/incremental-averageing
  r.N += 1
  r.Average = r.Average + (value - r.Average)/float32(r.N)

  // fmt.Printf("addValue: average[%f]\n", r.Average)
}

// I'm going to keep the data for the incremental averages in a map from action name
// to current average. Note that this is a static (not instance) variable. The requirement is for
// a library to keep track of averages, which suggests just a single set of averages.
var averages = make(map[string]*currentAverage)

// To establish a critical section around averages
var mutex = &sync.Mutex{}

// Accepts a JSON string of the following form:
//    {"action":"jump", "time":100}
// and maintains an average time for each action.
// If the input JSON cannot be parsed this throws an error.
func AddAction(s string) error {
  var time ActionTime
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
    averages[*time.Action] = new(currentAverage)
  }

  averages[*time.Action].addValue(*time.Time)
  mutex.Unlock()

  // fmt.Printf("AddAction: key[%s] value[%f]\n", *time.Action, *time.Time)

  return nil // no error
}

// Returns an empty string if there was an error.
func GetStats() string {
  var results = []ActionTime{}

  mutex.Lock()
  for action, average := range averages {
    //fmt.Printf("key[%s] value[%f]\n", action, average.Average)

    // Passing &(action) directly into the ActionTime struct results in a bug. I think it's a string
    // allocation issue. Instead, do it this way:
    var name = action

    next := ActionTime{&name, &average.Average}
    results = append(results, next)
  }
  mutex.Unlock()

  bytes, err := json.Marshal(results)
  if err != nil {
    return ""
  }

  return string(bytes)
}

// This was not in the requirements but was needed for testing. I think it also
// might be useful in a practical application.
func Reset() {
  mutex.Lock()
  averages = make(map[string]*currentAverage)
  mutex.Unlock()
}

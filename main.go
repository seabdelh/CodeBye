package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sedki-abdelhakim/chatbot"
	// Autoload environment variables in .env

	_ "github.com/joho/godotenv/autoload"
)

/////////////////////////////////
type UserStruct struct {
	Status string
	Result []struct {
		LastName                string
		Country                 string
		LastOnlineTimeSeconds   int
		City                    string
		Rating                  int
		FriendOfCount           int
		TitlePhoto              string
		Handle                  string
		Avatar                  string
		FirstName               string
		Contribution            int
		Organization            string
		Rank                    string
		MaxRating               int
		RegistrationTimeSeconds int
		MaxRank                 string
	}
}

type ContestStatus struct {
	Status string
	Result []struct {
		ID                  int
		ContestID           int
		CreationTimeSeconds int
		RelativeTimeSeconds int64
		Problem             struct {
			ContestID int
			Index     string
			Name      string
			Type      string
			Points    float64
			Tags      []string
		}
		Author struct {
			ContestID int
			Members   []struct {
				Handle string
			}
			ParticipantType  string
			Ghost            bool
			StartTimeSeconds int
		}
		ProgrammingLanguage string
		Verdict             string
		Testset             string
		PassedTestCount     int
		TimeConsumedMillis  int
		MemoryConsumedBytes int
	}
}

type ContestStandings struct {
	Status string
	Result struct {
		Contest struct {
			ID                  int
			Name                string
			Type                string
			Phase               string
			Frozen              bool
			DurationSeconds     int
			StartTimeSeconds    int
			RelativeTimeSeconds int
		}
		Problems []struct {
			ContestID int
			Index     string
			Name      string
			Type      string
			Points    float64
			Tags      []string
		}
		Rows []struct {
			Party struct {
				ContestID int
				Members   []struct {
					Handle string
				}
				ParticipantType  string
				Ghost            bool
				Room             int
				StartTimeSeconds int
			}
			Rank                  int
			Points                float64
			Penalty               int
			SuccessfulHackCount   int
			UnsuccessfulHackCount int
			ProblemResults        []struct {
				Points                    float64
				RejectedAttemptCount      int
				Type                      string
				BestSubmissionTimeSeconds int
			}
		}
	}
}

type Tags struct {
	Status string
	Result struct {
		Problems []struct {
			ContestID int
			Index     string
			Name      string
			Type      string
			Points    float64
			Tags      []string
		}
		ProblemStatistics []struct {
			ContestID   string
			Index       string
			SolvedCount int
		}
	}
}

/////////////////////////////////////////////

func chatbotProcess(session *chatbot.Session, message string) (string, error) {
	switch session.State {
	case 0:
		return handle0Out(session, message), nil
	case 1:
		return handle1Out(session, message), nil
	case 2:
		return handle2Out(session, message), nil
	case 3:
		return handle1In(session, message), nil

	}

	return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil
}
func handle0Out(session *chatbot.Session, message string) string {

	if validateHandle(message) {
		return handle1In(session, message)
	}
	return handle0In(session, message)
}
func handle0In(session *chatbot.Session, message string) string {
	session.state = 0
	return "Wrong handle, please enter a valid handle"

}

func handle1Out(session *chatbot.Session, message string) string {
	messageArr := strings.Split(message, " ")
	keyword := strings.ToLower(messageArr[0])
	handle := messageArr[1]
	problem := messageArr[3]
	var messageReply string

	switch keyword {
	case "did":
		if validateHandle(handle) && validateProblem(problem) {
			messageReply = handle4In(message)
		}
		break
	case "could":
		if validtag(messageArr[5]) {
			messageReply = handle2In(session, message)
		} else {
			messageReply = handle1In(session, message)
		}
		break
	}
	return messageReply
}

func handle1In(session *chatbot.Session, message string) string {
	session.State = 1
	return "So, how could I help you?"

}

func handle2In(session *chatbot.Session, message string) string {
	session.State = 2
	messageLower := strings.ToLower(message)
	ss := strings.Split(messageLower, " ")
	session.Tag = ss[5]
	return "What level"

}

func handle2Out(session *chatbot.Session, message string) string {
	if message == "easy" || message == "medium" || message == "hard" {
		return handle3In(session, message)
	}
	return handle2In(session, message)
}
func handle3In(session *chatbot.Session, msg string) string {
	session.State = 3
	// easy > 3000
	// 1000 < medium <3000
	// hard < 1000
	resp, _ := http.Get("http://codeforces.com/api/problemset.problems?tags=" + session.Tag)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	t := Tags{}
	json.Unmarshal(body, &t)
	prob := t.Result.ProblemStatistics
	l := len(prob)
	for i := 0; i < l; i++ {
		if msg == "easy" && prob[i].SolvedCount > 3000 {
			return "http://codeforces.com/problemset/problem/" + prob[i].ContestID + "/" + prob[i].Index
		} else if msg == "hard" && prob[i].SolvedCount < 1000 {
			return "http://codeforces.com/problemset/problem/" + prob[i].ContestID + "/" + prob[i].Index
		} else if msg == "medium" && prob[i].SolvedCount >= 1000 && prob[i].SolvedCount < 3000 {
			return "http://codeforces.com/problemset/problem/" + prob[i].ContestID + "/" + prob[i].Index
		}
	}
	return "sorry can't find a suitable problem"
}

func handle4In(message string) string {

	messageArr := strings.Split(message, " ")
	handle := messageArr[1]
	problem := messageArr[3]
	resp, _ := http.Get("http://codeforces.com/api/contest.status?contestId=" + problem[:len(problem)-1] + "&handle=" + handle)
	body, _ := ioutil.ReadAll(resp.Body)
	contestStatus := ContestStatus{}
	json.Unmarshal(body, &contestStatus)
	for _, result := range contestStatus.Result {
		if strings.EqualFold(result.Problem.Index, string(problem[len(problem)-1])) && result.Verdict == "OK" {
			return handle + " has solved problem: " + problem + " in " +
				strconv.Itoa(result.TimeConsumedMillis) + " milli seconds"
		}
	}
	return handle + " has not solved problem: " + problem
}
func validtag(tag string) bool {
	resp, _ := http.Get("http://codeforces.com/api/problemset.problems?tags=" + tag)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	t := Tags{}
	json.Unmarshal(body, &t)
	return len(t.Result.Problems) > 0
}
func validateHandle(handle string) bool {
	resp, _ := http.Get("http://codeforces.com/api/user.info?handles=" + handle)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	user := UserStruct{}
	json.Unmarshal(body, &user)
	return user.Status == "OK"

}
func validateProblem(problem string) bool {
	contestID := problem[:len(problem)-1]
	resp, _ := http.Get("http://codeforces.com/api/contest.standings?contestId=" + contestID + "&from=1&count=1")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	contest := ContestStandings{}
	json.Unmarshal(body, &contest)

	for _, value := range contest.Result.Problems {
		if strings.EqualFold(value.Index, string(problem[len(problem)-1])) { //do not change messageLower[0] to messageLower :)
			return true //problemID is valid
		}
	}

	return false //problemID is invalid

}
func main() {
	// Uncomment the following lines to customize the chatbot
	chatbot.WelcomeMessage = "Hi Mr. coder.Let's have fun, what is your codeforces handle ?"
	chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(chatbot.Engage(":" + port))
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sedki-abdelhakim/chatbot"
	// Autoload environment variables in .env

	_ "github.com/joho/godotenv/autoload"
)

// UserStruct respresents a user object returned from Codeforces API
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

// ContestStatus represents status of a contest and its submissions in Codeforces.
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

// ContestStandings respresents result of a contest and its problems in Codeforces.
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

// Tags struct respresents list of problems with a specific tag on Codeforces.
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

// ProjectObject represents a response of todo.ly when creating objects
type ProjectObject struct {
	Id 			string
	Content 	string
	ItemsCount 	string
	Icon 		string
	ItemType	string
	ParentId	string
	Collapsed	string
	ItemOrder	string
	Children	string
	IsProjectShared	string
	IsShareApproved	string
	IsOwnProject	string
	LastSyncedDateTime	string
	LastUpdatedDate		string
	Deleted		string
}

const (
	passLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randomPass(l int) string {
	p := make([]byte, l)
	for i := range p {
		p[i] = passLetters[rand.Intn(len(passLetters))]
	}
	return string(p)
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
	// cases 3,4,5,6 returning to state 1 with success
	case 3:
		return handle1Out(session, message), nil
	case 4:
		return handle1Out(session, message), nil
	case 5:
		return handle1Out(session, message), nil

	}

	return "", nil
}
func handle0Out(session *chatbot.Session, message string) string {

	if validateHandle(message) {
		session.Handel = message
		return handle1In(session, message)
	}
	return handle0In(session, message)
}
func handle0In(session *chatbot.Session, message string) string {
	session.State = 0
	return "Wrong handle, please enter a valid handle"

}

func handle1Out(session *chatbot.Session, message string) string {
	messageArr := strings.Split(message, " ")
	keyword := strings.ToLower(messageArr[0])
	var messageReply string

	switch keyword {
	case "did":
		if len(messageArr) > 3 && validateHandle(messageArr[1]) && validateProblem(messageArr[3]) {
			messageReply = handle4In(session, message)
		} else {
			messageReply = handle1In(session, message)
		}
		break
	case "could":
		if len(messageArr) > 5 && validtag(strings.ToLower(messageArr[5])) {
			messageReply = handle2In(session, message)
		} else {
			messageReply = handle1In(session, message)
		}
		break
	case "give":
		if len(messageArr) > 5 {
			messageReply = handle5In(session, messageArr[5])
		} else {
			messageReply = handle1In(session, message)

		}
		break
	case "coach":
		handle6In(session, message)
		break
	}
	if messageReply == "" {
		messageReply = "I did not get that !"
	}
	return messageReply
}

func handle1In(session *chatbot.Session, message string) string {
	messageRes := "So, how could I help you?"
	if session.State == 1 || session.State == 2 {
		messageRes = "I did not get that !"
	}

	session.State = 1
	return messageRes
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

func handle4In(session *chatbot.Session, message string) string {
	//about if someone solved a specific problem
	//format: did [handle] solved [problemID]
	session.State = 4

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

func handle5In(session *chatbot.Session, handle string) string {
	//about the progress of someone
	//format: give me some info about [handle]
	session.State = 5
	if strings.EqualFold("me", handle) {
		handle = session.Handel //it is handle not handel but who cares L :)
	}
	resp, _ := http.Get("http://codeforces.com/api/user.info?handles=" + handle)
	body, _ := ioutil.ReadAll(resp.Body)
	user := UserStruct{}
	json.Unmarshal(body, &user)
	if user.Status != "OK" {
		return "Sorry I could not get you the info you want"
	}
	return `First name: ` + user.Result[0].FirstName + ", " +
		`Last name: ` + user.Result[0].LastName + ", " +
		`Rating: ` + strconv.Itoa(user.Result[0].Rating) + ", " +
		`Country: ` + user.Result[0].Country + ", " +
		`Last online: ` + strconv.Itoa(user.Result[0].LastOnlineTimeSeconds) + ", " +
		`City: ` + user.Result[0].City + ", " +
		`Number of friends: ` + strconv.Itoa(user.Result[0].FriendOfCount) + ", " +
		`Title photo: ` + user.Result[0].TitlePhoto + ", " +
		`Handle: ` + user.Result[0].Handle + ", " +
		`Avatar: ` + user.Result[0].Avatar + ", " +
		`Contribution: ` + strconv.Itoa(user.Result[0].Contribution) + ", " +
		`Organization: ` + user.Result[0].Organization + ", " +
		`Rank: ` + user.Result[0].Rank + ", " +
		`Max rating: ` + strconv.Itoa(user.Result[0].MaxRating) + ", " +
		`Registration time: ` + strconv.Itoa(user.Result[0].RegistrationTimeSeconds) + ", " +
		`Max rank: ` + user.Result[0].MaxRank

}

func handle6In(session *chatbot.Session, message string) string {
	session.State = 6

	r := randomPass(4)
	generatedEmail := session.Handel + r + "@codebye.me"
	generatedPass := randomPass(8)
	fullName := session.Handel

	createTodoUser(generatedEmail, generatedPass, fullName)

	problems := getProblems()

	projectName := "CodeBye Plan " + r
	projectId := createTodoProject(projectName, generatedEmail, generatedPass)



}

func createTodoProject(projectName string, generatedEmail string, generatedPass string) string {
	XMLCreateProject := "<ProjectObject>" +
		"<Content>" + projectName + "</Content> " +
		"<Icon>4</Icon> " +
		"</ProjectObject>"

	client := &http.Client{}
    req, err := http.NewRequest("POST", "https://todo.ly/api/projects.xml", strings.NewReader(XMLCreateProject))
    req.SetBasicAuth(generatedEmail, generatedPass)
    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }
    bodyText, err := ioutil.ReadAll(resp.Body)


    s := string(bodyText)
    return s
}

func getProblems() []string {

	var tags = [...]string{"implementation", "dp", "geometry", "math", "greedy", "strings", "graphs", "trees", "games", "probabilities", "bitmasks", "combinatorics"}
	resp, _ := http.Get("http://codeforces.com/api/problemset.problems?tags=" + tags[rand.Intn(len(tags))])
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	t := Tags{}
	json.Unmarshal(body, &t)
	prob := t.Result.ProblemStatistics

	var problemURLs []string
	for i in range 4 {
		problemURLs = append(problemURLs, "http://codeforces.com/problemset/problem/"+prob[i].ContestID+"/"+prob[1].Index)
	}
	
	return problemURLs

}

func createTodoUser(generatedEmail string, generatedPass string, fullName string) {
	XMLCreateUser := "<UserObject>" +
		"<Email>" + generatedEmail + "</Email> " +
		"<FullName>" + fullName + "</FullName> " +
		"<Password>" + generatedPass + "</Password> " +
		"</UserObject>"

	resp, err := http.Post("https://todo.ly/api/user.xml", "text/xml", strings.NewReader(XMLCreateUser))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
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

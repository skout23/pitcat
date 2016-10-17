package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("admin-directory_v1-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {

	var given, family, email, password, better string

	given = "Purple"
	family = "Drank"
	email = "pdrink@skoutrocks.org"
	password = "StupidStupidStupidPass"

	better = "SCRUBBEDLIKEABOSS"

	// set variables based on input

	name := &admin.UserName{
		GivenName:  given,
		FamilyName: family,
	}
	user := &admin.User{
		Name:                      name,
		Password:                  password,
		PrimaryEmail:              email,
		ChangePasswordAtNextLogin: false,
	}

	dork := &admin.User{
		Password:     better,
		PrimaryEmail: email,
	}

	fmt.Println(user)

	imapsettings := &gmail.ImapSettings{
		Enabled: true,
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/admin-directory_v1-go-quickstart.json
	config, err := google.ConfigFromJSON(b, admin.AdminDirectoryUserScope, gmail.GmailSettingsBasicScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := admin.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}
	gsrv, err := gmail.New(client)

	// delete pdrink

	// srv.Users.Delete(email).Do()

	// insert Purple Drank user with bad password

	/*
		user2, err := srv.Users.Insert(user).Do()
		if err != nil {
			log.Fatalf("Cannot create user in domain. %v", err)
		} else {
			log.Printf("Succeed to create user: %v", user2)
		}
	*/

	// Grab current list of users Object r
	/*
		r, err := srv.Users.List().Customer("my_customer").MaxResults(10).
			OrderBy("email").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve users in domain.", err)
		}

		// iterate over r

		if len(r.Users) == 0 {
			fmt.Print("No users found.\n")
		} else {
			fmt.Print("Users:\n")
			for _, u := range r.Users {
				fmt.Printf("%s (%s)\tAdmin:()\n", u.PrimaryEmail, u.Name.FullName)
			}
		}
	*/

	// Patch pdrink user with better password

	srv.Users.Patch(email, dork).Do()
	ires, bop := gsrv.Users.Settings.UpdateImap(email, imapsettings).Do()
	gres, gbop := gsrv.Users.Settings.GetImap(email).Do()
	fmt.Println(ires)
	fmt.Println(bop)
	fmt.Println(gres)
	fmt.Println(gbop)

}

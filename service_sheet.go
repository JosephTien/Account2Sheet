package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "google.golang.org/api/sheets/v4"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)
type Sheet struct{
    rows [][]interface {}
    updated bool//sync problem!!!
    gid string
}
var(
    keypath = "key"+string(os.PathSeparator)
    defaultSpreadsheetId = "1zvYlacc1ESyAcBoxuOyLlZ_Uiilz5MA8b21_p_NzWng"
    srv *sheets.Service
    spreadsheets = make(map[string]*Sheet)//chach // manual edit may ruin the chach
)
/////////////////////////////////////////////////////////////////////////////////
func readList(spreadsheetId string, tableName string)bool{
    readRange := fmt.Sprintf("%s!A2:C",tableName)
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet: %v", err)
        return false
    }
    if len(resp.Values) == 0 {
        fmt.Println("No data found.")
    } else {
        for _, row := range resp.Values {
            hashstr := fmt.Sprintf("%s/%s",row[0],row[1])
            _, exist := spreadsheets[hashstr]
            if !exist{
                if len(row)==3{
                    spreadsheets[hashstr] = &Sheet{nil, false, fmt.Sprintf("%s",row[2])}
                }else{
                    spreadsheets[hashstr] = &Sheet{nil, false, ""}
                }
                
            }
        }
    }
    return true
}
/////////////////////////////////////////////////////////////////////////////////
func readInfos(spreadsheetId string, tableName string)bool{
    hashstr := fmt.Sprintf("%s/%s",spreadsheetId,tableName)
    _, exist := spreadsheets[hashstr]
    if exist{
        if spreadsheets[hashstr].updated{
            //return false //maybe manual edit, so don't trust the chach
        }
    }else{
        spreadsheets[hashstr] = &Sheet{nil, false, ""}
    }
    //-------------------------------------
    readRange := fmt.Sprintf("%s!A2:G",tableName)
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet: %v", err)
        return true
    }
    if len(resp.Values) == 0 {
        fmt.Println("No data found.")
    } else {
        
    }
    spreadsheets[hashstr].rows = resp.Values
    spreadsheets[hashstr].updated = true;
    return false
}
func addInfo(data Data)bool{
    spreadsheetId := data.SpreadsheetId
    tableName := data.TableName
    readInfos(spreadsheetId, tableName)
    hashstr := fmt.Sprintf("%s/%s", spreadsheetId, tableName)
    //-------------------------------------
    tarrow:=len(spreadsheets[hashstr].rows)+2
    str_calculate:="=G2-H2"
    str_wallet:="=IF(D2=\"Y\",G2-H2,0.00)"
    if tarrow>2{
        str_calculate=fmt.Sprintf("=I%d+G%d-H%d", tarrow-1, tarrow, tarrow)
        str_wallet=fmt.Sprintf("=IF(D%d=\"Y\",J%d+G%d-H%d,J%d)", tarrow, tarrow-1, tarrow, tarrow, tarrow-1)
    }
    writeRange := fmt.Sprintf("%s!A2",tableName)
    var vr sheets.ValueRange
    myval := []interface{}{data.Date, data.Item, data.State, data.Payer, data.Receipt, data.Reimburse, data.Income, data.Outcome, str_calculate, str_wallet}
    vr.Values = append(vr.Values, myval)
    _, err := srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet. %v", err)
        return true
    }
    //-------------------------------------
    spreadsheets[hashstr].rows = append(spreadsheets[hashstr].rows, myval)
    return false
}
/////////////////////////////////////////////////////////////////////////////////
func prePareSheetsService(){
    b, err := ioutil.ReadFile(keypath+"credentials.json")
    if err != nil {
            log.Fatalf("Unable to read client secret file: %v", err)
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
    if err != nil {
            log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client := getClient(config)
    srv, err = sheets.New(client)
    if err != nil {
            log.Fatalf("Unable to retrieve Sheets client: %v", err)
    }
    return
}

func getClient(config *oauth2.Config) *http.Client {
    // The file token.json stores the user's access and refresh tokens, and is
    // created automatically when the authorization flow completes for the first
    // time.
    tokFile := keypath+"token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
            tok = getTokenFromWeb(config)
            saveToken(tokFile, tok)
    }
    return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
            "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
            log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
            log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
            return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
            log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

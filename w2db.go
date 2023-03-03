// ================================================================================================
// File Name	: w2db.go
// Project		: w2db
// Author		: Holger Scheller
// Version		:
// Last Update	: 23.02.2023
// Description	: Backend functions for the w2ui library Version 1.5.
// ================================================================================================
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const nInsRec = 10000000

// customer record for w2grid
type customerRec struct {
	Recid   int    `json:"recid"`
	Usr     string `json:"usr"`
	Pwd     string `json:"pwd"`
	Title   string `json:"title"`
	Fname   string `json:"fname"`
	Lname   string `json:"lname"`
	Company string `json:"company"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZIP     string `json:"zip"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
}

// ================================================================================================
// Function		: w2grid
// Description	: request / response w2grid url option
// ================================================================================================
func w2grid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // set header to prevent blocked by CORS policy.
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	type jsonReq struct {
		Action string `json:"action"`
	}
	var jReq jsonReq
	sRes := ""
	db, err := sql.Open("sqlite3", "w2db.db")
	if err != nil {
		log.Printf("ERROR: %v", err)
		sRes = fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error())
		io.WriteString(w, sRes)
		return
	}
	defer db.Close()
	sReq := r.URL.Query().Get("request")
	err = json.Unmarshal([]byte(sReq), &jReq)
	switch jReq.Action {
	case "":
		sRes, err = gridData(db, sReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			sRes = fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error())
		}
		break
	case "delete":
		sRes, err = gridDelete(db, sReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			sRes = fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error())
		}
		break
	case "save":
		sRes, err = gridSave(db, sReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			sRes = fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error())
		}
		break
	}
	io.WriteString(w, sRes)
}

// ================================================================================================
// Function		: w2form
// Description	: request / response w2form url option
// ================================================================================================
func w2form(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // set header to prevent blocked by CORS policy.
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	type jsonReq struct {
		Cmd    string          `json:"cmd"`
		Recid  int             `json:"recid"`
		Name   string          `json:"name"`
		Record json.RawMessage `json:"record"`
	}
	var jReq jsonReq
	rec := customerRec{}
	db, err := sql.Open("sqlite3", "w2db.db")
	if err != nil {
		log.Printf("ERROR: %v", err)
		io.WriteString(w, fmt.Sprintf(`{"status":"error", "message": "%s"}`, err))
		return
	}
	defer db.Close()
	bBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: %v", err)
		io.WriteString(w, fmt.Sprintf(`{"status":"error", "message": "%s"}`, err))
		return
	}
	sBody, _ := url.QueryUnescape(string(bBody))
	sBody = strings.Replace(sBody, "request=", "", 1)
	err = json.Unmarshal([]byte(sBody), &jReq)

	if jReq.Cmd == "get" {
		if jReq.Recid == nInsRec {
			io.WriteString(w, fmt.Sprintf(`{"status":"success","record":{"recid":%d,"usr":"","pwd":"","title":"","fname":"","lname":"","company":"","street":"","city":"","state":"","zip":"","country":"","phone":""}}`, nInsRec))
			return
		}
		stmt, err := db.Prepare("SELECT [recid],[usr],[pwd],[title],[fname],[lname],[company],[street],[city],[state],[zip],[country],[phone] FROM [customer] WHERE [recid]=?")
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(jReq.Recid).Scan(&rec.Recid, &rec.Usr, &rec.Pwd, &rec.Title, &rec.Fname, &rec.Lname, &rec.Company, &rec.Street, &rec.City, &rec.State, &rec.ZIP, &rec.Country, &rec.Phone)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
			return
		}
		sRes, err := json.Marshal(rec)
		if err != nil {
			io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
			log.Printf("ERROR: %v", err)
			return
		}
		io.WriteString(w, fmt.Sprintf(`{"status":"success","record":%s}`, sRes))
		return
	} else if jReq.Recid > 0 {
		err = json.Unmarshal(jReq.Record, &rec)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()))
			return
		}
		if jReq.Recid == nInsRec {
			stmt, err := db.Prepare("INSERT INTO [customer] ([usr], [pwd], [title], [fname], [lname], [company], [street], [city], [state], [zip], [country], [phone]) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
				return
			}
			defer stmt.Close() // should be closed after use.
			if _, err := stmt.Exec(rec.Usr, rec.Pwd, rec.Title, rec.Fname, rec.Lname, rec.Company, rec.Street, rec.City, rec.State, rec.ZIP, rec.Country, rec.Phone); err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
				return
			}

		} else {
			stmt, err := db.Prepare("UPDATE [customer] SET [usr]=?, [pwd]=?, [title]=?, [fname]=?, [lname]=?, [company]=?, [street]=?, [city]=?, [state]=?, [zip]=?, [country]=?, [phone]=? WHERE [recid]=?")
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
				return
			}
			defer stmt.Close() // should be closed after use.
			if _, err := stmt.Exec(rec.Usr, rec.Pwd, rec.Title, rec.Fname, rec.Lname, rec.Company, rec.Street, rec.City, rec.State, rec.ZIP, rec.Country, rec.Phone, jReq.Recid); err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, fmt.Sprintf(`{"status":"error", "message":"%s"}`, err.Error()))
				return
			}
		}
		io.WriteString(w, `{"status":"success"}`)
		return
	}
}

// ================================================================================================
// Function		: gridData
// Description	: Create a json string with customer record's
// ================================================================================================
func gridData(db *sql.DB, sReq string) (string, error) {
	type jsonReq struct {
		Limit       int    `json:"limit"`
		Offset      int    `json:"offset"`
		SearchLogic string `json:"searchLogic"`
		Search      []struct {
			Field    string `json:"field"`
			Type     string `json:"type"`
			Operator string `json:"operator"`
			Value    string `json:"value"`
		} `json:"search"`
		Sort []struct {
			Field     string `json:"field"`
			Direction string `json:"direction"`
		} `json:"sort"`
	}
	var jReq jsonReq
	var aRec []*customerRec
	sSort := ""
	sSearch := ""
	nRec := 0
	sRes := ""
	err := json.Unmarshal([]byte(sReq), &jReq)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
	}
	if len(jReq.Sort) > 0 {
		sSort = "ORDER BY "
		for i := 0; i < len(jReq.Sort); i++ {
			sSort += fmt.Sprintf("[%s] %s,", jReq.Sort[i].Field, jReq.Sort[i].Direction)
		}
		sSort = sSort[:len(sSort)-1]
	}
	if len(jReq.Search) > 0 {
		for i := 0; i < len(jReq.Search); i++ {
			if sSearch == "" {
				sSearch = "WHERE ("
			} else {
				sSearch += fmt.Sprintf(" %s (", jReq.SearchLogic)
			}
			switch jReq.Search[i].Operator {
			case "is":
				sSearch += fmt.Sprintf("[%s] = '%s'", jReq.Search[i].Field, jReq.Search[i].Value)
			case "begins":
				sSearch += fmt.Sprintf("[%s] LIKE '%s'", jReq.Search[i].Field, jReq.Search[i].Value+"%")
			case "contains":
				sSearch += fmt.Sprintf("[%s] LIKE '%s'", jReq.Search[i].Field, "%"+jReq.Search[i].Value+"%")
			case "ends":
				sSearch += fmt.Sprintf("[%s] LIKE '%s'", jReq.Search[i].Field, "%"+jReq.Search[i].Value)
			case "before":
				sSearch += fmt.Sprintf("[%s] < '%s'", jReq.Search[i].Field, jReq.Search[i].Value)
			case "less":
				sSearch += fmt.Sprintf("[%s] < '%s'", jReq.Search[i].Field, jReq.Search[i].Value)
			case "after":
				sSearch += fmt.Sprintf("[%s] > '%s'", jReq.Search[i].Field, jReq.Search[i].Value)
			case "more":
				sSearch += fmt.Sprintf("[%s] > %s", jReq.Search[i].Field, jReq.Search[i].Value)
			case "between":
				sValue := jReq.Search[i].Value[:1]
				sValue = sValue[:len(sValue)-1]
				aValue := strings.Split(sValue, ",")
				sSearch += fmt.Sprintf("[%s] BETWEEN '%s' AND '%s'", jReq.Search[i].Field, aValue[0], aValue[1])
			}
			sSearch += ")"
		}
	} else {
		row := db.QueryRow("SELECT COUNT(1) AS [nRec] FROM [customer]")
		err = row.Scan(&nRec)
		if err != nil {
			nRec = -1
		}
	}
	sSql := fmt.Sprintf("SELECT [recid],[usr],[pwd],[title],[fname],[lname],[company],[street],[city],[state],[zip],[country],[phone] FROM [customer] %s %s LIMIT %d OFFSET %d", sSearch, sSort, jReq.Limit, jReq.Offset)
	rows, err := db.Query(sSql)
	if err != nil {
		log.Printf("ERROR: %v SQL: %s", err, sSql)
		return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
	}
	for rows.Next() {
		rec := new(customerRec)
		err := rows.Scan(&rec.Recid, &rec.Usr, &rec.Pwd, &rec.Title, &rec.Fname, &rec.Lname, &rec.Company, &rec.Street, &rec.City, &rec.State, &rec.ZIP, &rec.Country, &rec.Phone)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
		}
		aRec = append(aRec, rec)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: %v", err)
		return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
	}
	jRes, _ := json.Marshal(aRec)
	sRes = string(jRes)
	sRes = fmt.Sprintf(`{"status":"success","total":%d,"records":%s}`, nRec, sRes)
	return sRes, nil
}

// ================================================================================================
// Function		: gridData
// Description	: Save changes of w2grid inline edit record's. (w2grid button 'Save')
// ================================================================================================
func gridSave(db *sql.DB, sReq string) (string, error) {
	type jsonReq struct {
		Action  string            `json:"action"`
		Changes []json.RawMessage `json:"changes"`
	}
	var jReq jsonReq
	var sSql [3]string
	nRec := 0
	err := json.Unmarshal([]byte(sReq), &jReq)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
	}
	// create and execute an SQL-update query from the json fields
	for i := 0; i < len(jReq.Changes); i++ {
		var jField map[string]interface{}
		err := json.Unmarshal([]byte(jReq.Changes[i]), &jField)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
		}
		sSql[1] = ""
		sSql[2] = ""
		for key, value := range jField {
			if key == "recid" {
				sSql[1] = "[recid]=" + fmt.Sprint(value)
			} else {
				sSql[2] += `[` + key + `]="` + strings.Replace(fmt.Sprint(value), "\"", "'", -1) + `",`
			}
		}
		sSql[2] = sSql[2][:len(sSql[2])-1]
		sSql[0] = fmt.Sprintf("UPDATE [customer] SET %s WHERE %s", sSql[2], sSql[1])
		if _, err := db.Exec(sSql[0]); err != nil {
			log.Printf("ERROR: %v", err)
			return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
		}
		nRec += 1
	}
	return fmt.Sprintf(`{"status" : "success", "message" : "%d records updated."}`, nRec), nil
}

// ================================================================================================
// Function		: gridData
// Description	: Delete w2grid record's (w2grid button 'Delete')
// ================================================================================================
func gridDelete(db *sql.DB, sReq string) (string, error) {
	type jsonReq struct {
		Action string `json:"action"`
		Recid  []int  `json:"recid"`
	}
	var jReq jsonReq
	err := json.Unmarshal([]byte(sReq), &jReq)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
	}
	for i := 0; i < len(jReq.Recid); i++ {
		stmt, err := db.Prepare("DELETE FROM [customer] WHERE [recid]=?")
		if err != nil {
			log.Printf("ERROR: %v", err)
			return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(jReq.Recid[i]); err != nil {
			log.Printf("ERROR: %v", err)
			return fmt.Sprintf(`{"status":"error", "message": "%s"}`, err.Error()), err
		}
	}
	return `{"status" : "success"}`, nil
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	valueNum = 6
)

var (
	regex = regexp.MustCompile(`^ *(\S+) (\d+) (\S+) (\S{2,3}) (\S{2}) (\d{4}-\d{2}-\d{2}) *$`)
)

type fileConfig struct {
	In  string
	Out string
}

func parseFile(path string) int {
	invitations := make([]*invitation, 0)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	br := bufio.NewReader(f)
	for i := 0; ; i++ {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		invitations = append(invitations, genInvitation(line))
	}

	creates(invitations)

	return len(invitations)
}

func exportFile(path string) int {
	var buf bytes.Buffer
	buf.WriteString(
		"INSERT INTO invitations (name,money,subject,type,category,at,created_at,updated_at,deleted_at) VALUES\n",
	)

	values := findAll()
	total := len(values)

	for i := range values {
		deletedAt := "NULL"
		if values[i].DeletedAt != nil {
			deletedAt = fmt.Sprintf("'%s'", time2Str(*values[i].DeletedAt))
		}

		sql := fmt.Sprintf(
			"('%s',%d,'%s','%s','%s','%s','%s','%s',%s),\n",
			values[i].Name, values[i].Money, values[i].Subject, values[i].Type, values[i].Category,
			time2Str(values[i].At), time2Str(values[i].CreatedAt), time2Str(values[i].UpdatedAt), deletedAt,
		)

		if i == total-1 {
			sql = sql[:len(sql)-2] + ";"
		}

		buf.WriteString(sql)
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0666); err != nil {
		panic(err)
	}

	return total
}

func genInvitation(line []byte) *invitation {
	if !regex.Match(line) {
		panic(fmt.Errorf("该行与正则不匹配: %s", string(line)))
	}

	values := strings.Fields(string(line))
	if len(values) != valueNum {
		panic(fmt.Errorf("值个数不对[%d], %s", len(values), line))
	}

	return newInvitation(values)
}

func time2Str(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

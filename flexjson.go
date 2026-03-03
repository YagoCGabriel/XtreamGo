package main

import (
    "encoding/json"
    "strconv"
)

// FlexString aceita tanto string quanto number no JSON.
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
    // Tenta string primeiro
    var s string
    if err := json.Unmarshal(data, &s); err == nil {
        *f = FlexString(s)
        return nil
    }
    // Tenta número (int ou float)
    var n json.Number
    if err := json.Unmarshal(data, &n); err == nil {
        *f = FlexString(n.String())
        return nil
    }
    *f = ""
    return nil
}

func (f FlexString) String() string { return string(f) }

// FlexInt aceita tanto int quanto string no JSON.
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
    // Tenta número direto
    var n int
    if err := json.Unmarshal(data, &n); err == nil {
        *f = FlexInt(n)
        return nil
    }
    // Tenta string contendo número
    var s string
    if err := json.Unmarshal(data, &s); err == nil {
        n, err := strconv.Atoi(s)
        if err == nil {
            *f = FlexInt(n)
            return nil
        }
    }
    *f = 0
    return nil
}

func (f FlexInt) Int() int { return int(f) }
package prices

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func cmp(filldata map[string]interface{}, filled []PriceVol, t *testing.T, where string){

	if len(filled) != len(filldata){
		t.Error("Len mismatch on filled data", where,
			filldata,
			filled)
	}
	for _, got := range filled{
		pricestr := strconv.FormatFloat(got.Price, 'f', -1, 64)
		volstr := strconv.FormatFloat(got.Vol, 'f', -1, 64)
		entry, ok := filldata[pricestr]
		if !ok{
			t.Error("Failed to find price ", pricestr, where)
		}
		if volstr != entry{
			t.Error("Mismatch volumne for price ", pricestr, where)
		}
	}
}

func map2Triplets(m map[string]interface{})[]interface{}{
	ret := make([]interface{}, 3 * len(m))
	offset := 0
	for price, vol := range m{
		ret[offset + PRICE_INDEX] = price
		ret[offset + VOL_INDEX] = vol
		ret[offset + TIME_INDEX] = strconv.FormatInt(time.Now().Unix(), 10)
		offset += 3
	}
	return ret
}

func Test_FillAndUpdate(t *testing.T){
	filldata := map[string]interface{}{
		"6.72389": "120.34",
		"999.989": "10000.54332",
		"32.67723": "444.32",
	}
	positions := NewPositions()
	positions.Fill(filldata)
	fmt.Println("UPDATES ", positions.updated)

	updates := positions.GetUpdatesUnordered()
	if len(updates) > 0{
		t.Error("Updates set after fill", updates)
	}

	filled := positions.GetAllUnordered()
	cmp(filldata, filled, t, "first fill")

	updatedata := map[string]interface{}{
		"6.72389": "121.34",
		"32.67723": "0",
		"64.911211":"0.00012",
	}
	updatetriplets := map2Triplets(updatedata)
	positions.Update(updatetriplets)
	updated := positions.GetUpdatesUnordered()
	cmp(updatedata, updated, t, "first update")

	////updates should now be cleared
	updated = positions.GetUpdatesUnordered()
	if len(updated) > 0{
		t.Error("Updates not cleared after get ", updated)
	}

	////Finally test that the full data is correct - including the removal of zero vol
	expectedfinal := make(map[string]interface{})
	for price, vol := range filldata{
		uptvol, ok := updatedata[price]
		if ok && uptvol == "0"{
			continue
		}
		expectedfinal[price] = vol
	}
	for price, vol := range updatedata{
		if vol != "0" {
			expectedfinal[price] = vol
		}
	}
	filled = positions.GetAllUnordered()
	cmp(expectedfinal, filled, t, "final")
}


func cmpStrFloat(exp map[string]float64, posdata []PriceVol, t *testing.T){

	for pricestr, vol := range exp{
		expprice := toFloat(pricestr)
		found := false
		gotvol := float64(0)
		for _, got := range posdata{
			if got.Price == expprice{
				found = true
				gotvol = got.Vol
				break
			}
		}
		if !found{
			t.Error("Failed to find price ", expprice)
		}
		if gotvol != vol{
			t.Error("Vol mismatch at price ", expprice, gotvol, vol)
		}
	}
}
func Test_FillAndUpdateByStrFloat(t *testing.T){
	filldata := map[string]float64{
		"233.4566": 2330.56,
		"103.432": 998778,
		"6722.0921": 55429.3324,
	}
	pos := NewPositions()
	pos.FillStrFloat(filldata)
	posdata := pos.GetAllUnordered()
	cmpStrFloat(filldata, posdata, t)
	upddata := pos.GetUpdatesUnordered()
	if len(upddata) > 0{
		t.Error("Got update data after a fill")
	}

	update := map[string]float64{
		"103.432": 0,
		"6722.0921": 55429.7762,
		"6801.95": 42021,
	}
	///update the fill data map to the new expected
	for price, vol := range update {
		if vol != 0{
			filldata[price] = vol
		}else {
			delete(filldata, price)
		}
	}
	pos.UpdateStrFloat(update)
	upddata = pos.GetUpdatesUnordered()
	cmpStrFloat(update, upddata, t)
	////check update data is now empty
	upddata = pos.GetUpdatesUnordered()
	if len(upddata) > 0{
		t.Error("Updates still set after get")
	}

	///check the overall mapping is correct
	posdata = pos.GetAllUnordered()
	cmpStrFloat(filldata, posdata, t)


}
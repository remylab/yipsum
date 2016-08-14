package common

import (
    //"fmt"
    "testing"

    "github.com/stretchr/testify/assert"
)


func TestGetSentences(t *testing.T) {

    res := GetSentences("Speaking as a Providence resident: Pawtucket begins where the comfortable-walking-distance-to-the-Providence-station ends. Pawtucket and Central Falls are languishing. Providence, generally speaking, is more affordable than Boston to live in (which drives a lot of commuter traffic to Boston), but the presence of Brown is steadily driving rents up in the suburb closest to the Providence station. Pawtucket and Central Falls should be the affordable residential suburbs of Providence, but they're not (or, at least, not so much as they should be). Commuter rail to these towns would help revitalize them and ease some of Providence's present growing pains. The thrust of commuter rail into southern Rhode Island has failed because it has not responded to the needs of residents. In Wickford, where I grew up, very few residents take advantage of the weekday commuter lines running from the new Wickford Junction Station, since very few residents there commute to Boston. There are, however, plenty of retired folk who would love to take a weekend trip to Boston via train, but are stymied by the station being closed on weekends. If weekday commuter rail is to ever succeed in southern Rhode Island, it will not be in the short term.")

    assert.Equal(t, 10, len(res))

}

package seal

import (
	"testing"
	"strconv"
	"github.com/stretchr/testify/assert"
)

func TestSDB(t *testing.T){
	sdb:=NewSDB()
	str:="a secret data"
	b:=[]byte(str)
	label:="secret"
	assert.Equal(t,str,string(b))
	sdb.Insert(label,b)
	assert.NotEqual(t,str,string(b))
	sdb.Unlock(label)
	assert.Equal(t,str,string(b))
	sdb.Lock(label)
	assert.NotEqual(t,str,string(b))
	bb:=sdb.Get(label)
	assert.Equal(t,string(b),string(bb))
	sdb.Unlock(label)
	assert.Equal(t,str,string(bb))
	sdb.Lock(label)
	sdb.Delete(label)
	assert.Nil(t,sdb.Get(label))
	str2:=str
	label2:=label
	for i:=2; i<24; i++{
		str2=str2+"."
		b:=[]byte(str2)
		label2=label+strconv.Itoa(i)
		sdb.Insert(label2,b)
		assert.NotEqual(t,str2,string(b))
		sdb.Unlock(label2)
		assert.Equal(t,str2,string(b))
		sdb.Lock(label2)
		assert.NotEqual(t,str2,string(b))
	}
	

}

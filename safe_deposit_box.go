package seal

import (
	"crypto/rand"
	"sync"
	"encoding/binary"
	"github.com/spaolacci/murmur3"
	"errors"
)


type tkey struct {
	dataCount uint64
	key []byte
}
type tdata struct {
	locked bool
	data []byte
}

type SafeDepositBox struct {
	mut sync.Mutex
	keys map[uint32]*tkey
	data map[string]*tdata
}

func NewSDB()*SafeDepositBox{
	return &SafeDepositBox{
		keys:make(map[uint32]*tkey),
		data:make(map[string]*tdata),
	}
}


func obfuscate(k, dt []byte){

	for i:=0; i<len(dt); i+=8 {
		if i+7>=len(dt) {
			if i+3<len(dt) {
				dt[i],dt[i+1],dt[i+2],dt[i+3]=dt[i+2],dt[i+3],dt[i+1],dt[i]
			}
			break
		}
		dt[i],dt[i+1],dt[i+2],dt[i+3],dt[i+4],dt[i+5],dt[i+6],dt[i+7]=dt[i+1],dt[i+4],dt[i+5],dt[i+2],dt[i+7],dt[i],dt[i+3],dt[i+6]
	}
	a:=binary.BigEndian.Uint32(k[:4])
	b:=binary.BigEndian.Uint32(k[4:8])
        c:=binary.BigEndian.Uint32(k[8:12])
        d:=binary.BigEndian.Uint32(k[12:])
	rk:=make([]byte,16)
	VeryRandOP()
	for i:=0; i<len(dt); i+=16 {
		a+=b
		b^=c
		c-=d
		d^=a

		a^=b
		b-=c
		c^=d
		d+=a

		a-=b
		b^=c
		c+=d
		d^=a

		binary.BigEndian.PutUint32(rk[:4],a)
		binary.BigEndian.PutUint32(rk[4:8],b)
		binary.BigEndian.PutUint32(rk[8:12],c)
		binary.BigEndian.PutUint32(rk[12:],d)
		l:=0;
		for j:=i; j<i+16 && j<len(dt); j++{
			dt[j]^=rk[l]
			l++
		}

	}
	RandOP(1)
}

func deobfuscate(k,dt []byte){
	a:=binary.BigEndian.Uint32(k[:4])
        b:=binary.BigEndian.Uint32(k[4:8])
        c:=binary.BigEndian.Uint32(k[8:12])
        d:=binary.BigEndian.Uint32(k[12:])
	rk:=make([]byte,16)
	VeryRandOP()
	for i:=0; i<len(dt); i+=16 {
                a+=b
                b^=c
                c-=d
                d^=a

                a^=b
                b-=c
                c^=d
                d+=a

                a-=b
                b^=c
                c+=d
                d^=a

                binary.BigEndian.PutUint32(rk[:4],a)
                binary.BigEndian.PutUint32(rk[4:8],b)
                binary.BigEndian.PutUint32(rk[8:12],c)
                binary.BigEndian.PutUint32(rk[12:],d)
		l:=0
                for j:=i; j<i+16 && j<len(dt); j++{
                        dt[j]^=rk[l]
			l++
                }

        }
	RandOP(1)

	for i:=0; i<len(dt); i+=8{
		if i+7>=len(dt) {
			if i+3<len(dt) {
				dt[i+2],dt[i+3],dt[i+1],dt[i]=dt[i],dt[i+1],dt[i+2],dt[i+3]
			}
			break
		}
		dt[i+1],dt[i+4],dt[i+5],dt[i+2],dt[i+7],dt[i],dt[i+3],dt[i+6]=dt[i],dt[i+1],dt[i+2],dt[i+3],dt[i+4],dt[i+5],dt[i+6],dt[i+7]
	}

}

func (sdb *SafeDepositBox) Insert(label string,data []byte)error{
	if label=="" {
		return errors.New("label can't empty")
	}
	if data==nil || len(data)==0 {
		return errors.New("data can't nil or empty")
	}
	sdb.mut.Lock()
	defer sdb.mut.Unlock()
	_,labelExists:=sdb.data[label]
	if labelExists {
		return errors.New("label already in use")
	}
	h:=murmur3.New32()
	h.Write([]byte(label))
	kIndex:=h.Sum32()
	key,kIndexExists:=sdb.keys[kIndex]
	if !kIndexExists {
		_key:=make([]byte,16)
		rand.Read(_key)
		sdb.keys[kIndex]=&tkey{
			dataCount:1,
			key:_key,
		}
		key=sdb.keys[kIndex]
	}else {
		key.dataCount++
	}
	obfuscate(key.key,data)
	sdb.data[label]=&tdata{
		locked:true,
		data:data,
	}
	return nil
}

func (sdb *SafeDepositBox) Lock(label string)error{
	if label=="" {
		return nil
	}
	sdb.mut.Lock()
	defer sdb.mut.Unlock()
	data,labelExists:=sdb.data[label]
	if !labelExists {
		return errors.New("label not exists")
	}
	h:=murmur3.New32()
	h.Write([]byte(label))
	kIndex:=h.Sum32()
	key,kIndexExists:=sdb.keys[kIndex]
	if !kIndexExists {
		return errors.New("key not exists")
	}
	if data.locked {
		return nil
	}
	obfuscate(key.key,data.data)
	data.locked=true

	return nil
}

func (sdb *SafeDepositBox) Unlock(label string)error{
	if label=="" {
		return nil
	}
	sdb.mut.Lock()
	defer sdb.mut.Unlock()
	data,labelExists:=sdb.data[label]
	if !labelExists {
		return errors.New("label not exists")
	}
	h:=murmur3.New32()
	h.Write([]byte(label))
	kIndex:=h.Sum32()
	key,kIndexExists:=sdb.keys[kIndex]
	if !kIndexExists {
		return errors.New("key not exists")
	}
	deobfuscate(key.key,data.data)
	data.locked=false
	return nil
}

func (sdb *SafeDepositBox) Get(label string)[]byte{
	if label=="" {
		return nil
	}
	sdb.mut.Lock()
	defer sdb.mut.Unlock()
	data,ok:=sdb.data[label]
	if ok {
		return data.data
	}
	return nil
}

func (sdb *SafeDepositBox) Delete(label string){
	if label=="" {
		return
	}
	sdb.mut.Lock()
	defer sdb.mut.Unlock()
	data,labelExists:=sdb.data[label]
	h:=murmur3.New32()
	h.Write([]byte(label))
	kIndex:=h.Sum32()
	key,kIndexExists:=sdb.keys[kIndex]
	if labelExists {
		obfuscate(key.key,data.data)
		delete(sdb.data,label)
	}
	if kIndexExists {
		if labelExists {
			key.dataCount--
		}
		if key.dataCount<=0 {
			delete(sdb.keys,kIndex)
		}
	}
}

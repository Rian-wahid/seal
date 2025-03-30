package seal

import (
	"crypto/rand"
	"encoding/binary"
)

func RandOP(i int){
	
	r:=make([]byte,16)
	rand.Read(r)
	
	a:=binary.BigEndian.Uint32(r[:4])
	b:=binary.BigEndian.Uint32(r[4:8])
	c:=binary.BigEndian.Uint32(r[8:12])
	d:=binary.BigEndian.Uint32(r[12:16])
	for j:=0; j<i; j++{
		switch(a%3){
			case 0:
				a+=b
				a^=c
				c*=d
				d-=b
				b^=a
				d+=c
				break
			case 1:
				b-=a
				b*=d
				a+=c
				c^=a
				break
			case 2:
				c*=d
				d^=a
				break
		}

		switch(b%3){
                case 0:
                        	a+=b
                        	a^=c
                        	c*=d
                        	d-=b
				b^=a
				d+=c
                        	break
                	case 1:
                        	b-=a
                        	b*=d
                        	a+=c
                        	c^=a
                        	break
                	case 2:
                        	c*=d
                        	d^=a
                        	break
        	}

		switch(c%3){
                	case 0:
                        	a+=b
                        	a^=c
                        	c*=d
                        	d-=b
				b^=a
				d+=c
                        	break
                	case 1:
                        	b-=a
                        	b*=d
                        	a+=c
                        	c^=a
                        	break
                	case 2:
                        	c*=d
                        	d^=a
                        	break
        	}
	
		switch(d%3){
                	case 0:
                        	a+=b
                        	a^=c
                        	c*=d
                        	d-=b
				b^=a
				d+=c
                        	break
                	case 1:
                        	b-=a
                        	b*=d
                        	a+=c
                        	c^=a
                        	break
                	case 2:
                        	c*=d
                        	d^=a
                	        break
        	}
	}
}

func VeryRandOP(){

	b:=[]byte{0}
	rand.Read(b)
	RandOP(1+int(b[0]))
}

/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Monday 01 June 2020 - 20:00:19
** @Filename:				bsonBigInt.go
**
** @Last modified by:		Tbouder
*******************************************************************************/

package bigson

import	(
	"errors"
	"math/big"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type BigInt big.Int

// Int retreive the big.Int value of a BigInt object
func (b *BigInt) Int() *big.Int {
	return (*big.Int)(b)
}

// New create a new BigInt from a big.Int
func New(value *big.Int) *BigInt {
	b := BigInt(*value)
	return &b
}
// NewInt create a new BigInt from a int
func NewInt(value int) *BigInt {
	bigValue := big.NewInt(int64(value))
	b := BigInt(*bigValue)
	return &b
}
// Sum returns the sum of old + new
func Sum(old, new *BigInt) *BigInt {
	result := big.NewInt(0).Add(old.Int(), new.Int())
	return New(result)
}
// Sub returns the sub of old - new
func Sub(old, new *BigInt) *BigInt {
	result := big.NewInt(0).Sub(old.Int(), new.Int())
	return New(result)
}

//UnmarshalText implements the text Unmarshal interface
func (b *BigInt) UnmarshalText(text []byte) (err error) {
	var bigInt = new(big.Int)
	err = bigInt.UnmarshalText(text)
	if err != nil {
		value := big.NewInt(0)
		*b = BigInt(*value)
		return err
	}

	*b = BigInt(*bigInt)
	return nil
}

//MarshalText implements the text marshal interface
func (b *BigInt) MarshalText() (text []byte, err error) {
	if (b.Int().String() == `<nil>`) {
		return []byte("0"), nil
	}
	return []byte(b.Int().String()), nil
}

//MarshalBSONValue implements the bson.ValueMarshaler interface
func (b *BigInt) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if (b.Int().String() == `<nil>`) {
		return bsontype.String, bsoncore.AppendString(nil, "0"), nil
	}
	return bsontype.String, bsoncore.AppendString(nil, b.Int().String()), nil
}

// UnmarshalBSONValue is an interface implemented that can unmarshal a BSON
// value representation of themselves.
func (b *BigInt) UnmarshalBSONValue(tpe bsontype.Type, data []byte) error {
	str, _, ok := bsoncore.ReadString(data)
	if !ok {
		return errors.New(`impossible to read data as string`)
	}

	var bigInt = new(big.Int)
	bigIntByte, ok := bigInt.SetString(str, 10)

	if !ok {
		bigIntByte = big.NewInt(0)
		*b = BigInt(*bigIntByte)
		return nil
	}
	*b = BigInt(*bigIntByte)
	return nil
}

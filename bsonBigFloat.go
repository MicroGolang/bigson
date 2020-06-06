/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Monday 01 June 2020 - 20:00:19
** @Filename:				bsonBigFloat.go
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

// BigFloat is a wrapper over big.Float to implement only unmarshalText
// for json decoding.
type BigFloat big.Float

// Float retreive the big.Float value of a BigFloat object
func (b *BigFloat) Float() *big.Float {
	return (*big.Float)(b)
}

// NewBigFloat create a new BigFloat from a big.Float
func NewBigFloat(value *big.Float) *BigFloat {
	b := BigFloat(*value)
	return &b
}
// NewFloat create a new BigFloat from a int
func NewFloat(value float64) *BigFloat {
	bigValue := big.NewFloat(value)
	b := BigFloat(*bigValue)
	return &b
}
// SetString create a new BigFloat from a string
func (b *BigFloat) SetString(value string) (*BigFloat, bool) {
	bigValue, ok := big.NewFloat(0).SetString(value)
	if (!ok) {
		return nil, false
	}
	bi := BigFloat(*bigValue)
	return &bi, true
}

// String return a representation of b as a string
func (b *BigFloat) String() string {
	return b.Float().String()
}

//UnmarshalText implements the text Unmarshal interface
func (b *BigFloat) UnmarshalText(text []byte) (err error) {
	var bigFloat = new(big.Float)
	err = bigFloat.UnmarshalText(text)
	if err != nil {
		value := big.NewFloat(0)
		*b = BigFloat(*value)
		return err
	}

	*b = BigFloat(*bigFloat)
	return nil
}

//MarshalText implements the text marshal interface
func (b *BigFloat) MarshalText() (text []byte, err error) {
	if (b.Float().Text('f', -1) == `<nil>`) {
		return []byte("0"), nil
	}
	return []byte(b.Float().Text('f', -1)), nil
}

//MarshalBSONValue implements the bson.ValueMarshaler interface
func (b *BigFloat) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if (b.Float().Text('f', -1) == `<nil>`) {
		return bsontype.String, bsoncore.AppendString(nil, "0"), nil
	}
	return bsontype.String, bsoncore.AppendString(nil, b.Float().Text('f', -1)), nil
}

// UnmarshalBSONValue is an interface implemented that can unmarshal a BSON
// value representation of themselves.
func (b *BigFloat) UnmarshalBSONValue(tpe bsontype.Type, data []byte) error {
	str, _, ok := bsoncore.ReadString(data)
	if !ok {
		return errors.New(`impossible to read data as string`)
	}

	var bigFloat = new(big.Float)
	bigFloatByte, ok := bigFloat.SetString(str)

	if !ok {
		bigFloatByte = big.NewFloat(0)
		*b = BigFloat(*bigFloatByte)
		return nil
	}
	*b = BigFloat(*bigFloatByte)
	return nil
}

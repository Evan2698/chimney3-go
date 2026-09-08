package core

import (
	"errors"
	"io"
	"log"
	"net"
)

func ReadXBytes(bytes uint32, buffer []byte, con net.Conn) ([]byte, error) {
	//defer utils.Trace("readXBytes.readXBytes")()
	if bytes == 0 {
		return nil, errors.New("0 bytes can not read! ")
	}
	if con == nil {
		return nil, errors.New("nil connection")
	}
	if uint64(bytes) > uint64(len(buffer)) {
		return nil, errors.New("buffer is too small")
	}

	var index uint32
	var err error
	var n int
	for {
		n, err = con.Read(buffer[index:])
		if err != nil {
			// only log actual errors (EOF will be surfaced to caller)
			log.Println("error on read_bytes_from_socket", n, err)
			break
		}
		if n == 0 {
			err = io.ErrNoProgress
			break
		}
		index = index + uint32(n)
		if index >= bytes {
			break
		}
	}
	if index == bytes {
		err = nil
	}

	// final result logged only on error for quieter normal operation
	if err != nil {
		log.Println("read result size:", index, err)
	}
	return buffer[:index], err
}

func WriteXBytes(buffer []byte, con net.Conn) (int, error) {
	//defer utils.Trace("writeXBytes.writeXBytes")()
	nbytes := uint32(len(buffer))
	if con == nil {
		return 0, errors.New("nil connection")
	}
	if nbytes == 0 {
		return 0, nil
	}
	var index uint32 = 0
	var err error
	var n int
	for {
		n, err = con.Write(buffer[index:])
		if err != nil {
			log.Println("write bytes error:", n, err)
			break
		}
		if n == 0 {
			err = io.ErrShortWrite
			break
		}
		index = index + uint32(n)
		if index >= nbytes {
			break
		}
	}
	if index == nbytes {
		err = nil
	}

	if err != nil {
		log.Println("writeXBytes error:", n, err)
	}

	return int(index), err
}

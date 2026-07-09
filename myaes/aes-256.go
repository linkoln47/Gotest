package main

const (
	blockSize       = 16
	keySize         = 32
	numRounds       = 14
	Nk              = 8 // 4, 6 или 8 слов
	numRoundKeys    = numRounds + 1
	expandedKeySize = blockSize * numRoundKeys // 240 bytes
)

func init() {
	// Шаг 1: мультипликативные обратные. Ищем перебором (init выполняется раз).
	var inv [256]byte
	for x := 1; x < 256; x++ {
		for y := 1; y < 256; y++ {
			if gmul(byte(x), byte(y)) == 1 {
				inv[x] = byte(y)
				break
			}
		}
	}
	// Шаг 2: аффинное преобразование s = b ^ (b<<<1) ^ (b<<<2) ^ (b<<<3) ^ (b<<<4) ^ 0x63
	for i := 0; i < 256; i++ {
		b := inv[i]
		s := b ^ rotl8(b, 1) ^ rotl8(b, 2) ^ rotl8(b, 3) ^ rotl8(b, 4) ^ 0x63
		sbox[i] = s
		invSbox[s] = byte(i)
	}
}

// ============================================================================
// 1. Арифметика в конечном поле GF(2^8)
// ----------------------------------------------------------------------------
// Байты в AES — это элементы поля GF(2^8). Сложение = XOR. Умножение — по
// модулю неприводимого многочлена AES: x^8 + x^4 + x^3 + x + 1 (= 0x11B).
// ============================================================================

// xtime умножает элемент поля на x (то есть на 0x02) с приведением по модулю.
// Сдвиг влево на 1 бит = умножение на x; если «вылез» бит x^8, вычитаем
// (XOR-им) неприводимый многочлен 0x1B (младшие 8 бит от 0x11B).
func xtime(value byte) byte {
	if value&0x80 != 0 {
		return (value << 1) ^ 0x1b
	}
	return value << 1
}

// gmul — умножение двух элементов GF(2^8) методом «русского крестьянина».
func gmul(a, b byte) byte {
	var p byte
	for i := 0; i < 8; i++ {
		if b&1 != 0 {
			p ^= a // прибавляем текущую степень, если бит установлен
		}
		b >>= 1
		a = xtime(a) // переходим к следующей степени x^k
	}
	return p
}

// ============================================================================
// 2. S-box и обратный S-box
// ----------------------------------------------------------------------------
// Вместо «магической» таблицы из 256 чисел генерируем её из первых принципов:
//   1) мультипликативный обратный элемент в GF(2^8) (для 0 — сам 0);
//   2) аффинное преобразование над битами.
// Это и есть источник нелинейности AES.
// ============================================================================

var sbox, invSbox [256]byte

// rotl8 циклически вращает байт влево на shift бит.
func rotl8(x byte, shift int) byte {
	return (x << shift) | (x >> (8 - shift))
}

// ============================================================================
// 3. Расширение ключа (Key Expansion / Key Schedule)
// ----------------------------------------------------------------------------
// Из исходного ключа получаем numRounds+1 раундовых ключей по 16 байт.
// Слово = 4 байта (uint32). Nk слов в ключе, Nb=4 слова в блоке.
// ============================================================================

func subWord(w uint32) uint32 {
	return uint32(sbox[w>>24])<<24 |
		uint32(sbox[(w>>16)&0xff])<<16 |
		uint32(sbox[(w>>8)&0xff])<<8 |
		uint32(sbox[w&0xff])
}

func rotWord(w uint32) uint32 { return (w << 8) | (w >> 24) }

// keyExpansion возвращает срез раундовых ключей и число раундов numRounds.
func keyExpansion(key [keySize]byte) [numRoundKeys][blockSize]byte {
	total := 4 * (numRounds + 1) // всего слов в расписании

	w := make([]uint32, total)
	for i := 0; i < Nk; i++ {
		w[i] = uint32(key[4*i])<<24 | uint32(key[4*i+1])<<16 |
			uint32(key[4*i+2])<<8 | uint32(key[4*i+3])
	}

	rc := byte(1) // Rcon: rc[1]=1, далее rc[i]=xtime(rc[i-1])
	for i := Nk; i < total; i++ {
		temp := w[i-1]
		switch {
		case i%Nk == 0:
			temp = subWord(rotWord(temp)) ^ (uint32(rc) << 24)
			rc = xtime(rc)
		case i%Nk == 4: // дополнительный SubWord только для AES-256
			temp = subWord(temp)
		}
		w[i] = w[i-Nk] ^ temp
	}

	// Раскладываем слова в 16-байтные раундовые ключи (column-major).
	var roundKeys [numRoundKeys][blockSize]byte
	for r := 0; r < numRoundKeys; r++ {
		for c := 0; c < 4; c++ {
			word := w[4*r+c]
			roundKeys[r][4*c+0] = byte(word >> 24)
			roundKeys[r][4*c+1] = byte(word >> 16)
			roundKeys[r][4*c+2] = byte(word >> 8)
			roundKeys[r][4*c+3] = byte(word)
		}
	}
	return roundKeys
}

// ============================================================================
// 4. Операции раунда над состоянием (state)
// ----------------------------------------------------------------------------
// Состояние — 16 байт. Логически это матрица 4x4, заполняемая по столбцам:
// байт с индексом i стоит в позиции (row=i%4, col=i/4). То есть
//   state[4*col + row] = элемент (row, col).
// ============================================================================

// SubBytes: нелинейная замена каждого байта через S-box.
func subBytes(s *[16]byte) {
	for i := range s {
		s[i] = sbox[s[i]]
	}
}
func invSubBytes(s *[16]byte) {
	for i := range s {
		s[i] = invSbox[s[i]]
	}
}

// ShiftRows: строку r циклически сдвигаем влево на r позиций.
// new(r,c) = old(r, (c+r) mod 4).
func shiftRows(s *[16]byte) {
	var t [16]byte
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			t[4*c+r] = s[4*((c+r)%4)+r]
		}
	}
	*s = t
}
func invShiftRows(s *[16]byte) {
	var t [16]byte
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			t[4*c+r] = s[4*((c-r+4)%4)+r]
		}
	}
	*s = t
}

// MixColumns: каждый столбец умножается на фиксированный многочлен в GF(2^8).
// Матрица: [2 3 1 1; 1 2 3 1; 1 1 2 3; 3 1 1 2].
func mixColumns(s *[16]byte) {
	for c := 0; c < 4; c++ {
		s0, s1, s2, s3 := s[4*c], s[4*c+1], s[4*c+2], s[4*c+3]
		s[4*c+0] = xtime(s0) ^ (xtime(s1) ^ s1) ^ s2 ^ s3 // 2·s0 ^ 3·s1 ^ s2 ^ s3
		s[4*c+1] = s0 ^ xtime(s1) ^ (xtime(s2) ^ s2) ^ s3
		s[4*c+2] = s0 ^ s1 ^ xtime(s2) ^ (xtime(s3) ^ s3)
		s[4*c+3] = (xtime(s0) ^ s0) ^ s1 ^ s2 ^ xtime(s3)
	}
}

// InvMixColumns: обратная матрица [14 11 13 9; 9 14 11 13; 13 9 14 11; 11 13 9 14].
func invMixColumns(s *[16]byte) {
	for c := 0; c < 4; c++ {
		s0, s1, s2, s3 := s[4*c], s[4*c+1], s[4*c+2], s[4*c+3]
		s[4*c+0] = gmul(14, s0) ^ gmul(11, s1) ^ gmul(13, s2) ^ gmul(9, s3)
		s[4*c+1] = gmul(9, s0) ^ gmul(14, s1) ^ gmul(11, s2) ^ gmul(13, s3)
		s[4*c+2] = gmul(13, s0) ^ gmul(9, s1) ^ gmul(14, s2) ^ gmul(11, s3)
		s[4*c+3] = gmul(11, s0) ^ gmul(13, s1) ^ gmul(9, s2) ^ gmul(14, s3)
	}
}

// AddRoundKey: побайтовый XOR состояния с раундовым ключом.
func addRoundKey(s *[16]byte, rk [16]byte) {
	for i := range s {
		s[i] ^= rk[i]
	}
}

// ============================================================================
// 5. Блочный шифр: Cipher / InvCipher
// ============================================================================

type Cipher struct {
	roundKeys [numRoundKeys][blockSize]byte
}

// New создаёт шифр по ключу 32 байта (AES-256).
func New(key [keySize]byte) *Cipher {
	return &Cipher{keyExpansion(key)}
}

// EncryptBlock шифрует ровно один 16-байтный блок.
func (c *Cipher) EncryptBlock(in [16]byte) [16]byte {
	s := in
	addRoundKey(&s, c.roundKeys[0]) // начальный whitening

	for round := 1; round < numRounds; round++ { // основные раунды
		subBytes(&s)
		shiftRows(&s)
		mixColumns(&s)
		addRoundKey(&s, c.roundKeys[round])
	}

	subBytes(&s) // последний раунд — без MixColumns
	shiftRows(&s)
	addRoundKey(&s, c.roundKeys[numRounds])
	return s
}

// DecryptBlock — обратный порядок операций.
func (c *Cipher) DecryptBlock(in [16]byte) [16]byte {
	s := in
	addRoundKey(&s, c.roundKeys[numRounds])

	for round := numRounds - 1; round >= 1; round-- {
		invShiftRows(&s)
		invSubBytes(&s)
		addRoundKey(&s, c.roundKeys[round])
		invMixColumns(&s)
	}

	invShiftRows(&s)
	invSubBytes(&s)
	addRoundKey(&s, c.roundKeys[0])
	return s
}

// ============================================================================
// 6. Дополнение PKCS#7
// ----------------------------------------------------------------------------
// Добавляем N байт со значением N, где N = сколько не хватает до кратности
// блоку (1..16). Если данные уже кратны — добавляем целый блок из 16 байт 0x10,
// чтобы снятие дополнения было однозначным.
// ============================================================================

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	n := len(data)
	pad := int(data[n-1])
	return data[:n-pad], nil
}

// ============================================================================
// 7. Режимы работы: ECB и CBC
// ============================================================================

// --- ECB: каждый блок шифруется независимо. НЕБЕЗОПАСЕН: одинаковый открытый
//     текст даёт одинаковый шифр (виден рисунок данных). Только для учёбы. ---

func (c *Cipher) EncryptECB(plaintext []byte) []byte {
	data := pkcs7Pad(plaintext, 16)
	out := make([]byte, len(data))
	for i := 0; i < len(data); i += 16 {
		var b [16]byte
		copy(b[:], data[i:i+16])
		enc := c.EncryptBlock(b)
		copy(out[i:], enc[:])
	}
	return out
}

func (c *Cipher) DecryptECB(ciphertext []byte) ([]byte, error) {
	out := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += 16 {
		var b [16]byte
		copy(b[:], ciphertext[i:i+16])
		dec := c.DecryptBlock(b)
		copy(out[i:], dec[:])
	}
	return pkcs7Unpad(out)
}

// --- CBC: перед шифрованием каждый блок XOR-ится с предыдущим шифрблоком
//     (первый — с IV). IV должен быть случайным и уникальным на каждое
//     сообщение (16 байт). IV не секретен, но не должен повторяться. ---

func (c *Cipher) EncryptCBC(plaintext, iv []byte) ([]byte, error) {
	data := pkcs7Pad(plaintext, 16)
	out := make([]byte, len(data))
	prev := make([]byte, 16)
	copy(prev, iv)
	for i := 0; i < len(data); i += 16 {
		var b [16]byte
		for j := 0; j < 16; j++ {
			b[j] = data[i+j] ^ prev[j] // XOR с предыдущим шифрблоком/IV
		}
		enc := c.EncryptBlock(b)
		copy(out[i:], enc[:])
		copy(prev, enc[:])
	}
	return out, nil
}

func (c *Cipher) DecryptCBC(ciphertext, iv []byte) ([]byte, error) {
	out := make([]byte, len(ciphertext))
	prev := make([]byte, 16)
	copy(prev, iv)
	for i := 0; i < len(ciphertext); i += 16 {
		var b [16]byte
		copy(b[:], ciphertext[i:i+16])
		dec := c.DecryptBlock(b)
		for j := 0; j < 16; j++ {
			out[i+j] = dec[j] ^ prev[j]
		}
		copy(prev, ciphertext[i:i+16])
	}
	return pkcs7Unpad(out)
}

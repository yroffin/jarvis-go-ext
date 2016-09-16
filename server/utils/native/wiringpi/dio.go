package wiringpi

// #cgo arm CFLAGS: -marm
// #cgo arm LDFLAGS: -lwiringPi
// extern void scheduler_realtime();
// extern void scheduler_standard();
import "C"

const (
	OUTPUT = 1
	HIGH   = 1
	LOW    = 0
)

// On : send on
func On(pin int, sender uint64, interruptor uint64) {
	Send(pin, sender, interruptor, HIGH)
}

// Off : send off
func Off(pin int, sender uint64, interruptor uint64) {
	Send(pin, sender, interruptor, LOW)
}

// Calcul le nombre 2^chiffre indiqué, fonction utilisé par itob pour la conversion decimal/binaire
func power2(power int) uint64 {
	var integer uint64
	integer = 1
	for i := 0; i < power; i++ {
		integer *= 2
	}
	return integer
}

// Convertis un nombre en binaire, nécessite le nombre, et le nombre de bits souhaité en sortie (ici 26)
// Stocke le résultat dans le tableau global "bit2"
func itob(integer uint64, length int) []int {
	var bit2 = make([]int, length, length)
	for i := 0; i < length; i++ {
		if (integer / power2(length-1-i)) == 1 {
			integer -= power2(length - 1 - i)
			bit2[i] = 1
		} else {
			bit2[i] = 0
		}
	}
	return bit2
}

// Envois d'une pulsation (passage de l'etat haut a l'etat bas)
// 1 = 310µs haut puis 1340µs bas
// 0 = 310µs haut puis 310µs bas
func SendBit(pin int, b bool) {
	if b {
		DigitalWrite(pin, HIGH)
		DelayMicroseconds(310) //275 orinally, but tweaked.
		DigitalWrite(pin, LOW)
		DelayMicroseconds(1340) //1225 orinally, but tweaked.
	} else {
		DigitalWrite(pin, HIGH)
		DelayMicroseconds(310) //275 orinally, but tweaked.
		DigitalWrite(pin, LOW)
		DelayMicroseconds(310) //275 orinally, but tweaked.
	}
}

// Envoie d'une paire de pulsation radio qui definissent 1 bit réel : 0 =01 et 1 =10
// c'est le codage de manchester qui necessite ce petit bouzin, ceci permet entre autres de dissocier les données des parasites
func SendPair(pin int, b int) {
	if b == 1 {
		SendBit(pin, true)
		SendBit(pin, false)
	} else {
		SendBit(pin, false)
		SendBit(pin, true)
	}
}

// Fonction d'envois du signal
// recoit en parametre un booleen définissant l'arret ou la marche du matos (true = on, false = off)
func transmit(pin int, value int, bit2 []int, bit2Interruptor []int) {
	// Sequence de verrou anoncant le départ du signal au recepeteur
	DigitalWrite(pin, HIGH)
	DelayMicroseconds(275) // un bit de bruit avant de commencer pour remettre les delais du recepteur a 0
	DigitalWrite(pin, LOW)
	DelayMicroseconds(9900) // premier verrou de 9900µs
	DigitalWrite(pin, HIGH) // high again
	DelayMicroseconds(275)  // attente de 275µs entre les deux verrous
	DigitalWrite(pin, LOW)  // second verrou de 2675µs
	DelayMicroseconds(2675)
	DigitalWrite(pin, HIGH) // On reviens en état haut pour bien couper les verrous des données

	// Envoie du code emetteur (272946 = 1000010101000110010  en binaire)
	for i := 0; i < 26; i++ {
		SendPair(pin, bit2[i])
	}

	// Envoie du bit définissant si c'est une commande de groupe ou non (26em bit)
	SendPair(pin, 0)

	// Envoie du bit définissant si c'est allumé ou eteint 27em bit)
	SendPair(pin, value)

	// Envoie des 4 derniers bits, qui représentent le code interrupteur, ici 0 (encode sur 4 bit donc 0000)
	// nb: sur  les télécommandes officielle chacon, les interrupteurs sont logiquement nommés de 0 à x
	// interrupteur 1 = 0 (donc 0000) , interrupteur 2 = 1 (1000) , interrupteur 3 = 2 (0100) etc...
	for i := 0; i < 4; i++ {
		if bit2Interruptor[i] == 0 {
			SendPair(pin, 0)
		} else {
			SendPair(pin, 1)
		}
	}

	DigitalWrite(pin, HIGH) // coupure données, verrou
	DelayMicroseconds(275)  // attendre 275µs
	DigitalWrite(pin, LOW)  // verrou 2 de 2675µs pour signaler la fermeture du signal
}

// Send : send flow control to receptor
func Send(pin int, sender uint64, interruptor uint64, value int) {
	var bit2 []int
	var bit2Interruptor []int

	C.scheduler_realtime()

	PinMode(pin, OUTPUT)

	// 26 bit Identifiant emetteur
	// convertion du code de l'emetteur (ici 8217034) en code binaire
	bit2 = itob(sender, 26)
	bit2Interruptor = itob(interruptor, 4)

	for i := 0; i < 5; i++ {
		transmit(pin, value, bit2, bit2Interruptor) // envoyer ON
		Delay(10)                                   // attendre 10 ms (sinon le socket nous ignore)
	}

	C.scheduler_standard()
}

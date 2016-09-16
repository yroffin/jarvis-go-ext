#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>
#include <sys/time.h>
#include <time.h>
#include <stdlib.h>
#include <sched.h>

int wiringPiSetupInit() {
	if (setuid(0)) {
		perror("setuid");
		return 1;
	}

	//Si on ne trouve pas la librairie wiringPI, on arrête l'execution
	if (wiringPiSetup() == -1) {
		perror("Librairie Wiring PI introuvable, veuillez lier cette librairie...");
		return -1;
	}
}

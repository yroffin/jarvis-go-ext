package cmd

/*
#include <stdio.h>
#include <string.h>
#include<signal.h>
#include<unistd.h>
#include <stdlib.h>
void sig_handler(int signo)
{
  if (signo == SIGHUP)
    fprintf(stderr,"received SIGHUP and continue\n");
}
void  systemFork() {
    int child = fork();
    if(child > 0) {
        fprintf(stderr, "Fork daemon pid = %d\n", child);
        exit(0);
    } else {
        fprintf(stderr, "Fork\n");
        // Child process
        if (signal(SIGHUP, sig_handler) == SIG_ERR)
        return;
    }
}
*/
import "C"

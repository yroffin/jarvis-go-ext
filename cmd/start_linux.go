package cmd

/*
#include <stdio.h>
#include<signal.h>
#include<unistd.h>
void sig_handler(int signo)
{
  if (signo == SIGHUP)
    fprintf(stderr,"received SIGHUP and continue\n");
}
void  systemFork() {
    int child = fork();
    if(child == 0) {
        exit(0);
    } else {
        // Child process
        if (signal(SIGHUP, sig_handler) == SIG_ERR)
        return;
    }
}
*/
import "C"

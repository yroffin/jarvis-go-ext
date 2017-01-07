/**
 * Copyright 2017 Yannick Roffin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package cmd

/*
#include <stdio.h>
#include <string.h>
#include<signal.h>
#include<unistd.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>

void sig_handler(int signo)
{
  if (signo == SIGHUP)
    fprintf(stderr,"received SIGHUP and continue\n");
}

void  systemSighup() {
    if (signal(SIGHUP, sig_handler) == SIG_ERR) {
        fprintf(stderr, "Error, while setting signal handler\n");
    }
}

void  systemFork() {
    int child = fork();
    if(child > 0) {
        fprintf(stderr, "Fork daemon pid = %d\n", child);
        exit(0);
    } else {
        fprintf(stderr, "Fork\n");
        // Child process
        if (signal(SIGHUP, sig_handler) == SIG_ERR) {
            fprintf(stderr, "Error, while setting signal handler\n");
        }
        //unmask the file mode
        umask(0);
        //set new session
        setsid();
    }
}
*/
import "C"

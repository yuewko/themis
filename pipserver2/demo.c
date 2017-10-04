/************************************************************************/
/* demo.c - CFI demo.                                                   */
/* A simple non-performant demonstration on how to use basic            */
/* functionality in the RuleSpace CFI API.                              */
/* Copyright (c) 2011 Symantec Corporation. All rights reserved.        */
/* Use of this product is subject to license terms.                     */
/************************************************************************/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include "cfi.h"

#define M_LOCK(p)         pthread_mutex_lock(&p)
#define M_UNLOCK(p)       pthread_mutex_unlock(&p)
#define M_DESTROYLOCK(p)  pthread_mutex_destroy(&p)


#define THREADS 16
pthread_t thread[THREADS];

const char* cfg_fn = "conf.rd.direct";  /* CFI config filename, sane default */

pthread_mutex_t m_lock = PTHREAD_MUTEX_INITIALIZER;

/* Generic file loader. returns content and sets size of content
   returned in "len". If returned buffer is NULL, "len" is not
   set. filename can be NULL, wherein this function returns NULL. */
static void*
readfile(const char* filename, int* len)
{
  void* buf = NULL;
  FILE* file;

  if(filename) {
    file = fopen(filename, "r");
    if(file) {
      fseek(file, 0, SEEK_END);
      *len = ftell(file);
      fseek(file, 0, SEEK_SET);
      buf = malloc(*len);
      fread(buf, 1, *len, file);
      fclose(file);
    }
  }
  return buf;
}

int verbose = 0;

/* Basic wrapper around cfi_rate that demonstrates the lifecycle or
   CFI creation, rating, response cracking and cleanup. We simply
   print the result to stdout; Returns true if successful, false
   otherwise. */
rate(void* cfi,
     const char* url, int urlsz,
     const char* hdr, int hdrsz,
     const char* data, int datasz)
{
  void* rsp;              /* A rating response */
  char cat[CFI_CATMAX],   /* buffer to hold category names */
       txt[1024];         /* buffer for misc. text data    */
  void* tag = 0;          /* Tag used in cfi_rate */
  int val, sz = sizeof(txt), okay = 0;

  if(cfi) {
      /* Step 2: Now that we have a valid CFI handle,
         call cfi_rate to get a rating on the provided content. */
      if((rsp = cfi_rate(cfi, tag,
                         url, urlsz,
                         hdr, hdrsz,
                         data, datasz)) && verbose) {
        okay = 1;  
        /* Step 3: Collect rating data from the response    */

        /* If we looked up an URL, try display the found    */
        /* URL, also known as a FURL                        */
	fprintf(stdout, "URL: %.*s", urlsz, url);
        if (urlsz) {
          if ((cfi_get_furl(rsp, txt, &sz) && sz)) {
            fprintf(stdout, "FURL: %.*s\n", sz, txt);
          } else {
            fprintf(stdout, "FURL: <none>\n");
            return 0;
          }
        }

        switch(cfi_get_stat(rsp)) {
        case CFI_RSP_NOCAT:
          fprintf(stdout, "rate: NO CATEGORY\n");
          break;
        case CFI_RSP_UNKNOWN:
          fprintf(stdout, "rate: UNKNOWN\n");
          break;
        case CFI_RSP_RATED:
          cat[0] = 0;
          while(cfi_get_cat(rsp, cat, sizeof(cat), &val)) {
            fprintf(stdout, "rate: RATED, %s %d\n", cat, val);
          }
          break;
        case CFI_RSP_UNRATEABLE:
          fprintf(stdout, "rate: UNRATEABLE\n");
          okay = 0;
          break;
        case CFI_RSP_ERROR:
        default:
          sz = sizeof(txt);
          if (cfi_get_emsg(rsp, txt, &sz) && sz) {
            fprintf(stdout, "rate: ERROR: %s\n", txt);
          } else
            fprintf(stdout, "rate: ERROR\n");
          okay = 0;
          break;
        }
        /* Step 4: Clean-up, we are done rating */
        cfi_free_response(&rsp); /* Free the response; Note: we are
                                    passing the address of the response. */
      }

  }
  return okay;
}

pthread_mutex_t mutex1 = PTHREAD_MUTEX_INITIALIZER;

void *rate_loop()
{

    void * r_cfi = cfi_setup((char*) cfg_fn, stdout);
    
   if (!r_cfi){
      printf("NOTE: something went wrong with cfi_setup\n");
      return;
    }
    FILE* file;
    char *domain = NULL;
    size_t len = 0;
    ssize_t chars_read = 0;

    file = fopen("./test_domains", "r");
    if(file) {
      struct timeval t0, t1, t2;
      int count = 0;
      gettimeofday(&t0, NULL);
      int min = 100, max = 0;
      while(1){
         chars_read = getline(&domain, &len, file);
         if (chars_read < 1) {
            break;
         }
         gettimeofday(&t1, NULL);
         if (rate(r_cfi,
                domain, (domain) ? (int)strlen(domain) : 0,
                NULL, 0, NULL, 0)) {
         count++;
         gettimeofday(&t2, NULL);
         if ((min > (int)(t2.tv_usec-t1.tv_usec)) && ((int)(t2.tv_usec-t1.tv_usec> 0))){
            min = (int)t2.tv_usec-t1.tv_usec;
         }
         if (max < (int)(t2.tv_usec-t1.tv_usec)) {
            max = (int)(t2.tv_usec-t1.tv_usec);
         }
        }
      } /* while loop */

      gettimeofday(&t1, NULL);
      printf("\nNOTE: count %d", count);
      printf("   NOTE: %d min, %d max\n", min, max);
      printf("NOTE: rate all %ld sec.   ", t1.tv_sec-t0.tv_sec);
      printf("NOTE: %ld usec.\n", t1.tv_usec-t0.tv_usec);

      fclose(file);
    }

    cfi_end(&r_cfi);

}

void create_threads(void)
{
    int i;
    for(i =0; i < THREADS; i++) {
        printf("thread create loop\n");
        pthread_create(&thread[i], NULL, &rate_loop, NULL);
    }
}

/* add support for stdcall library */
#ifdef _WIN32
# define __CDECL __cdecl
#else
# define __CDECL
#endif

/* main that simply prepares and loads files and then calls the above
   rate function. */
int __CDECL
main(int argc, char **argv)
{
  int i, xval = 0;
  const char* url = NULL;     /* Opt URL to be rated */
  const char* hdr_fn = NULL;  /* Opt header filename (used in rating data) */
  const char* data_fn = NULL; /* Opt data filename */
  char *hdr, *data;           /* pointers to header and content data */
  int hdrsz, datasz;          /* ...and their respective sizes */
  void * r_cfi = NULL;

  for(i = 1; i < argc; i++) {
    if(!strcmp(argv[i], "-c")) {
      if(i < argc - 1) {
        cfg_fn = argv[++i];
      }
    } else if(!strcmp(argv[i], "-u")) {
      if(i < argc - 1) {
        url = argv[++i];
      }
    } else if(!strcmp(argv[i], "-h")) {
      if(i < argc - 1) {
        hdr_fn = argv[++i];
      }
    } else if(!strcmp(argv[i], "-d")) {
      if(i < argc - 1) {
        data_fn = argv[++i];
      }
    }
  }
  /* Print out disclaimer and legal-speak */
  printf(
    "************************************************************************\n"
    "* Throughput-limited Symantec Technology Demonstration                 *\n"
    "*                                                                      *\n"
    "* This sample program is intended to provide example usage of the      *\n"
    "* simplified CFI API and is not meant for production use. Limitations  *\n"
    "* include extremely limited throughput and synchronous behavior.       *\n");
  printf(
    "*                                                                      *\n"
    "* Contact support for additional information on advanced API usage.    *\n"
    "*                                                                      *\n"
    "* Copyright (c) 2011 Symantec Corporation. All rights reserved.        *\n"
    "************************************************************************\n");

    
    r_cfi = cfi_setup((char*) cfg_fn, stdout);
    cfi_list_opts_t list_opts;
    memset(&list_opts, 0, sizeof(list_opts));
    list_opts.lname = "list-name";
    list_opts.path = "location-of-the-RD-file";
    if (cfi_list_create_opts(r_cfi, &list_opts) == 1) {
        create_threads();
        for(i =0; i < THREADS; i++) {
            pthread_join(thread[i], NULL);
        }
    } else {
        fprintf(stdout, "Failed to open direct list wcd-r0.rd\n");
    }
    cfi_end(&r_cfi);
}

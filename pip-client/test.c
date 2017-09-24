#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <syslog.h>
#include <time.h>
#include "test.h"
#include "ts.h"

static TS_Handle TsHandle=NULL;
FILE *log_FP = NULL;


void my_log_function(TS_Log_Level level, TS_Log_Area areas, const char *message);


// return a string which is a list of category names
char* rate_url(TS_Handle ts_handle, const char *url, int verbosity)
{
     char *ts_result = NULL;    

     TS_Url parsed_url;
     TS_Attributes attributes;
     TS_Categories categories;

     char cat_names[4000];
     int num_cats = 0;
     //unsigned *cat_array=NULL;
     unsigned cat_array[100];

     char delimiter[] = ", ";
     int delimiter_len = strlen(delimiter);

     if (TS_OK != TS_AttributesCreate(ts_handle, &attributes)) {
	  printf("TS_AttributesCreate failed. Abort.\n");
	  return ts_result;
     }
     // printf("After TS_AttributesCreate()\n");

     if (TS_OK != TS_CategoriesCreate(ts_handle, &categories)) {
	  printf("TS_CategoriesCreate failed. Abort.\n");
	  TS_AttributesDestroy(ts_handle, &attributes);
	  return ts_result;
     }
     // printf("After TS_CategoriesCreate()\n");

     if (TS_OK != TS_CategoriesCategoryRemoveAll(ts_handle, categories)) {
	  printf("TS_CategoriesCategoryRemoveAll failed. Abort.\n");
	  TS_AttributesDestroy(ts_handle, &attributes);
	  TS_CategoriesDestroy(ts_handle, &categories);
	  //TS_HandleDestroy(&ts_handle);
	  return ts_result;
     }
     // printf("After TS_CategoriesCategoryRemoveAll()\n");

     if (TS_OK != TS_UrlCreate(ts_handle, &parsed_url)) {
	  printf("TS_UrlCreate failed. Abort.\n");
	  TS_AttributesDestroy(ts_handle, &attributes);
	  TS_CategoriesDestroy(ts_handle, &categories);
	  //TS_HandleDestroy(&ts_handle);
	  return ts_result;
     }
     // printf("After TS_UrlCreate()\n");
     //printf("parse url: %s\n", url);
     if (TS_OK != TS_UrlParse(ts_handle,
			      url,
			      NULL,
			      parsed_url)) {
	  printf("TS_UrlParse failed. Abort.\n");
	  TS_AttributesDestroy(ts_handle, &attributes);
      TS_CategoriesDestroy(ts_handle, &categories);
      TS_UrlDestroy(ts_handle, &parsed_url);
      return ts_result;
     }
     //printf("After TS_UrlParse()\n");

     if (TS_OK != TS_RateUrl(
	      ts_handle,
	      parsed_url,
	      attributes,
	      categories,
	      NULL,
	      0,
	      TS_CAT_SET_LOADED,
	      0,
	      NULL)) {
	  printf("TS_RateUrl failed. Abort.\n");
	  TS_AttributesDestroy(ts_handle, &attributes);
      TS_CategoriesDestroy(ts_handle, &categories);
      TS_UrlDestroy(ts_handle, &parsed_url);
      return ts_result;
     }

     // Get categories number
     
     if (TS_OK != TS_CategoriesCount(ts_handle, categories, &num_cats)) {
	  printf("Get categories number error!\n");
     }
     //printf("After TS_CategoriesCount()\n");
     //printf("num_cats: %d\n", num_cats);

     // Get categories codes
     // initialize the category array
     int i = 0;
     for (i=0; i<num_cats+1; i++) {
         cat_array[i] = 0;
     }
	 
     // ignore errors
     if (TS_OK != TS_CategoriesToArray(ts_handle, categories, cat_array, &num_cats) ) {
	  printf("Failed to get category code array!\n");
     }

     //printf("After TS_CategoriesToArray()\n");

#if 1
     int len = 0;
     len = sizeof(cat_names) - 1;
     if (TS_OK != TS_CategoriesToString(
	      ts_handle,
	      categories,
	      TS_LANGUAGE_ENGLISH,
	      TS_ENCODING_UTF8,
	      delimiter,
	      delimiter_len,
	      cat_names,
	      &len)) {
	  printf("TS_CategoriesToString failed. Abort.\n");
	  goto done;
     } else {
	  cat_names[len] = '\0';
	  if (strlen(cat_names) <= 1) {
	       // un-categoried:
           ts_result = strdup("uncategorized");
           
           if (verbosity >= 1)
    		    printf("x URL: '%s' is uncategorized!\n", url);
	  }
	  else {
           // categorized:
           char cat_codes[300];
           char tmpRet[4000];
           char code[10];
           int i;
           strcpy(cat_codes, "");
           for (i=0; i<num_cats; i++) {
			    //printf(" %u ", cat_array[i]);
                sprintf(code, " %u", cat_array[i]);
                strcat(cat_codes, code);
		   }
           //printf("cat codes: %s\n", cat_codes);
           
           sprintf(tmpRet, "%s %s", cat_names, cat_codes);
           // return the result
           ts_result = strdup(tmpRet);

	       if (verbosity >= 2) {
		    //
		    printf("URL: '%s' is categorized as :'%s'\n", url, ts_result);

		    printf("\n");
	       }
	  }
      //ts_result = strdup(cat_names);
     }
     // printf("After TS_CategoriesToString()\n");
#endif 
done:  
     TS_AttributesDestroy(ts_handle, &attributes);
     TS_CategoriesDestroy(ts_handle, &categories);
     TS_UrlDestroy(ts_handle, &parsed_url);
    // XXXXX: crashes here
    //  if (cat_array != NULL) {
	//      free(cat_array);
    //  }
    return ts_result;
}


void DestroySDK() {
    TS_HandleDestroy(&TsHandle);

}

int InitSDK() {
  //TS_Database_Access db_access_mode;
  const char *returned_serial = NULL;
  const char *errors = NULL;
  const char *client_cert = NULL;
  const char *client_key = NULL;
  const char *trustedsource_server_cert = NULL;

  // set log file
	log_FP = fopen("logfile", "a");
	if (NULL == log_FP)
	{
		log_FP = stdout;
	}

  
  if (TS_Init() != TS_OK)
  {
    fprintf(stderr, "TS_Init Failed. Abort.\n");
    return 0;
  }
  
  if (TS_OK != TS_HandleCreate(
	   &TsHandle,
	   "SF6S-HH37-G34G-X75H",
	   NULL,
	   "Infoblox",
	   "1"))
  {
       fprintf(stderr, "TS_HandleCreate failed. Abort.\n");
       return 0;
  }

  	/*
	 * Set the log level to info and log all areas.
	 */
	if (TS_OK != TS_LogLevelSet(
			TsHandle,
			TS_LOG_LEVEL_INFO,
			TS_LOG_AREA_ALL))
	{
		printf("TS_LogFunctionSet failed. Abort.\n");
		//TS_HandleDestroy(&TsHandle);
		return 0;
	}

	/*
	 * Set the log function to use my_log_function().
	 */
	if (TS_OK != TS_LogFunctionSet(TsHandle, my_log_function)) {
		printf("TS_LogFunctionSet failed. Abort.\n");
		//TS_HandleDestroy(&TsHandle);
		return 0;
	}


  if (TS_OK != TS_ActivateTrustedSource(
	   TsHandle,
	   TS_ACTIVATION_SERVER_DEFAULT,
	   NULL,
	   NULL,
	   &returned_serial,
	   &client_cert,
	   &client_key,
	   &trustedsource_server_cert,
	   &errors))
  {
       if (NULL == errors)
       {
	 fprintf(stderr, "Error during activation\n");
       }
       else
       {
	 fprintf(stderr, "Error from server: %s\n", errors);
       }
  }

  /*  db_access_mode = TS_DATABASE_ACCESS_MEMORY; */
  if (TS_OK != TS_DatabaseLoad( TsHandle,
				"data.db",
				TS_DATABASE_ACCESS_MEMORY,
				TS_CAT_SET_LATEST))
  {
       fprintf(stderr, "TS_DatabaseDownload failed. Abort!\n");
       //TS_HandleDestroy(&TsHandle);
       return 0;
  }
  else
  {
       printf("Local DB is loaded successfully.\n");       
  }

   return 0;
}

// Customized log function
void my_log_function(TS_Log_Level level, TS_Log_Area areas, const char *message) {

	if (NULL != log_FP) {
		//
		fprintf(log_FP, "%s\n", message);
		fflush(log_FP);
	}

	if ((TS_LOG_LEVEL_ERROR == level) && (TS_LOG_AREA_DATABASE_LOAD & areas)) {
		syslog(LOG_ERR, "%s\n", message);
	}
	return;
}


/*
 rate_url() returns a malloc'ed string that
 must be free'ed in calling code
*/
char *RateUrl(const char *url) {
    return rate_url(TsHandle, url, 0);
}

#if 0
// a test main
//  gcc test.c -Wall -g -lpthread -ldl -lrt -lssl -lcrypto libts.a -o client
int main(int argc, char *argv[]) {
     InitSDK();
     time_t start_time, current_time;
	 time(&start_time);
	 //printf("start_time : %s", start_time);
	 time(&current_time);
     char *ret_results = NULL;
     int duration = 5;
	 while (current_time < start_time + duration)
	 {
	    ret_results = RateUrl("little-porn0000.com");
        //ret_results = RateUrl("www.thesun.co.uk");
        //ret_results = RateUrl("123.uncategorized-domain.com");      

        if (ret_results != NULL) {
            printf("ret_result: %s\n", ret_results);
            free(ret_results);
        }

	 	time(&current_time);
	 }

     DestroySDK();

     //RateUrl("www.thesun.co.uk");
     return 0;
}

#endif


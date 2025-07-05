# Command Prompt
```bash
  jira issue view INGSVC-4929 --raw
```

# OUTPUT
{
  "expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations,customfield_10303.requestTypePractice",
  "id": "45305",
  "self": "https://thetalake.atlassian.net/rest/api/3/issue/45305",
  "key": "INGSVC-4929",
  "fields": {
    "parent": {
      "id": "42813",
      "key": "EPIC-2327",
      "self": "https://thetalake.atlassian.net/rest/api/3/issue/42813",
      "fields": {
        "summary": "Enhanced: Webex Meetings eComms Archive supports WSMP Private and Breakout Room messages ",
        "status": {
          "self": "https://thetalake.atlassian.net/rest/api/3/status/10136",
          "description": "Enhancement is fully documented and explained to all stakeholders, ready to ship (or is already shipped).",
          "iconUrl": "https://thetalake.atlassian.net/images/icons/statuses/generic.png",
          "name": "PM Customer Ready",
          "id": "10136",
          "statusCategory": {
            "self": "https://thetalake.atlassian.net/rest/api/3/statuscategory/3",
            "id": 3,
            "key": "done",
            "colorName": "green",
            "name": "Done"
          }
        },
        "priority": {
          "self": "https://thetalake.atlassian.net/rest/api/3/priority/3",
          "iconUrl": "https://thetalake.atlassian.net/images/icons/priorities/medium_new.svg",
          "name": "Medium",
          "id": "3"
        },
        "issuetype": {
          "self": "https://thetalake.atlassian.net/rest/api/3/issuetype/10000",
          "id": "10000",
          "description": "A big user story that needs to be broken down. Created by JIRA Software - do not edit or delete.",
          "iconUrl": "https://thetalake.atlassian.net/images/icons/issuetypes/epic.svg",
          "name": "Epic",
          "subtask": false,
          "hierarchyLevel": 1
        }
      }
    },
    "statusCategory": {
      "self": "https://thetalake.atlassian.net/rest/api/3/statuscategory/3",
      "id": 3,
      "key": "done",
      "colorName": "green",
      "name": "Done"
    },
    "resolution": {
      "self": "https://thetalake.atlassian.net/rest/api/3/resolution/10000",
      "id": "10000",
      "description": "Work has been completed on this issue.",
      "name": "Done"
    },
    "customfield_10510": null,
    "customfield_10506": null,
    "customfield_10748": "1_*:*_2_*:*_2435222086_*|*_3_*:*_2_*:*_699425488_*|*_5_*:*_1_*:*_48545_*|*_10001_*:*_1_*:*_87826693",
    "customfield_10507": null,
    "customfield_10508": null,
    "lastViewed": null,
    "labels": [],
    "aggregatetimeoriginalestimate": null,
    "issuelinks": [],
    "assignee": null,
    "components": [],
    "customfield_10841": null,
    "customfield_10842": null,
    "subtasks": [],
    "reporter": {
      "self": "https://thetalake.atlassian.net/rest/api/3/user?accountId=6080b9dcb8927500729602d7",
      "accountId": "6080b9dcb8927500729602d7",
      "emailAddress": "kannan@thetalake.com",
      "avatarUrls": {
        "48x48": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "24x24": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "16x16": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "32x32": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png"
      },
      "displayName": "Kannan Appachi",
      "active": true,
      "timeZone": "America/Los_Angeles",
      "accountType": "atlassian"
    },
    "customfield_10840": null,
    "customfield_10837": null,
    "customfield_10838": null,
    "progress": {
      "progress": 0,
      "total": 0
    },
    "customfield_10839": null,
    "votes": {
      "self": "https://thetalake.atlassian.net/rest/api/3/issue/INGSVC-4929/votes",
      "votes": 0,
      "hasVoted": false
    },
    "worklog": {
      "startAt": 0,
      "maxResults": 20,
      "total": 0,
      "worklogs": []
    },
    "issuetype": {
      "self": "https://thetalake.atlassian.net/rest/api/3/issuetype/10105",
      "id": "10105",
      "description": "A user story. Created by JIRA Software - do not edit or delete.",
      "iconUrl": "https://thetalake.atlassian.net/images/icons/issuetypes/story.svg",
      "name": "Story",
      "subtask": false,
      "hierarchyLevel": 0
    },
    "project": {
      "self": "https://thetalake.atlassian.net/rest/api/3/project/10000",
      "id": "10000",
      "key": "INGSVC",
      "name": "Ingestion Service",
      "projectTypeKey": "software",
      "simplified": false,
      "avatarUrls": {
        "48x48": "https://thetalake.atlassian.net/rest/api/3/universal_avatar/view/type/project/avatar/10203",
        "24x24": "https://thetalake.atlassian.net/rest/api/3/universal_avatar/view/type/project/avatar/10203?size=small",
        "16x16": "https://thetalake.atlassian.net/rest/api/3/universal_avatar/view/type/project/avatar/10203?size=xsmall",
        "32x32": "https://thetalake.atlassian.net/rest/api/3/universal_avatar/view/type/project/avatar/10203?size=medium"
      }
    },
    "customfield_10396": null,
    "customfield_11365": null,
    "resolutiondate": "2025-03-28T16:17:01.144-0700",
    "watches": {
      "self": "https://thetalake.atlassian.net/rest/api/3/issue/INGSVC-4929/watchers",
      "watchCount": 2,
      "isWatching": false
    },
    "customfield_11470": null,
    "customfield_11472": null,
    "customfield_11471": null,
    "customfield_11473": null,
    "customfield_11596": null,
    "customfield_11469": null,
    "customfield_11468": null,
    "customfield_10936": null,
    "updated": "2025-04-15T11:06:23.544-0700",
    "customfield_11580": null,
    "customfield_10370": null,
    "timeoriginalestimate": null,
    "customfield_11582": null,
    "customfield_10371": null,
    "description": {
      "type": "doc",
      "version": 1,
      "content": [
        {
          "type": "paragraph",
          "content": [
            {
              "type": "text",
              "text": "More Details in the Epic"
            }
          ]
        }
      ]
    },
    "customfield_10372": null,
    "customfield_11581": null,
    "customfield_10373": null,
    "customfield_11584": null,
    "customfield_11583": null,
    "customfield_10374": null,
    "customfield_11586": null,
    "customfield_11465": null,
    "customfield_10375": null,
    "customfield_11464": null,
    "customfield_11585": null,
    "customfield_11467": null,
    "timetracking": {},
    "customfield_11466": null,
    "customfield_11579": null,
    "customfield_10369": null,
    "customfield_11578": null,
    "customfield_10006": "EPIC-2327",
    "summary": "Webex Meetings eComms Archive : Implement Changes Required to Support WSMP Private and Breakout Room Messages",
    "customfield_11571": null,
    "customfield_10360": null,
    "customfield_11570": null,
    "customfield_11573": null,
    "customfield_10000": "{pullrequest={dataType=pullrequest, state=MERGED, stateCount=4}, json={\"cachedValue\":{\"errors\":[],\"summary\":{\"pullrequest\":{\"overall\":{\"count\":4,\"lastUpdated\":\"2025-03-28T10:45:12.568-0700\",\"stateCount\":4,\"state\":\"MERGED\",\"dataType\":\"pullrequest\",\"open\":false},\"byInstanceType\":{\"bitbucket\":{\"count\":4,\"name\":\"Bitbucket Cloud\"}}}}},\"isStale\":true}}",
    "customfield_11572": null,
    "customfield_11575": null,
    "customfield_10364": null,
    "customfield_10001": null,
    "customfield_11695": null,
    "customfield_10365": null,
    "customfield_10002": [],
    "customfield_11574": null,
    "customfield_11577": null,
    "customfield_10366": null,
    "customfield_11576": null,
    "customfield_10367": null,
    "customfield_11568": null,
    "customfield_10115": null,
    "customfield_10357": null,
    "customfield_11567": null,
    "customfield_10358": null,
    "environment": null,
    "customfield_10359": null,
    "customfield_11569": null,
    "duedate": null,
    "comment": {
      "comments": [
        {
          "self": "https://thetalake.atlassian.net/rest/api/3/issue/45305/comment/74547",
          "id": "74547",
          "author": {
            "self": "https://thetalake.atlassian.net/rest/api/3/user?accountId=614cd19ea995ad0073e155e9",
            "accountId": "614cd19ea995ad0073e155e9",
            "emailAddress": "aj@thetalake.com",
            "avatarUrls": {
              "48x48": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "24x24": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "16x16": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "32x32": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png"
            },
            "displayName": "Ardeshir Javaherchi",
            "active": true,
            "timeZone": "America/Los_Angeles",
            "accountType": "atlassian"
          },
          "body": {
            "type": "doc",
            "version": 1,
            "content": [
              {
                "type": "paragraph",
                "content": [
                  {
                    "type": "text",
                    "text": "Verified this on dev1 using datum 865931 which is the primary datum. Breakout room chats were captured in datums 865929 and 865930"
                  }
                ]
              }
            ]
          },
          "updateAuthor": {
            "self": "https://thetalake.atlassian.net/rest/api/3/user?accountId=614cd19ea995ad0073e155e9",
            "accountId": "614cd19ea995ad0073e155e9",
            "emailAddress": "aj@thetalake.com",
            "avatarUrls": {
              "48x48": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "24x24": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "16x16": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png",
              "32x32": "https://secure.gravatar.com/avatar/6d61d9cb98e6962d2865ba4fdb788da7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FAJ-2.png"
            },
            "displayName": "Ardeshir Javaherchi",
            "active": true,
            "timeZone": "America/Los_Angeles",
            "accountType": "atlassian"
          },
          "created": "2025-04-15T11:06:19.762-0700",
          "updated": "2025-04-15T11:06:19.762-0700",
          "jsdPublic": true
        }
      ],
      "self": "https://thetalake.atlassian.net/rest/api/3/issue/45305/comment",
      "maxResults": 1,
      "total": 1,
      "startAt": 0
    },
    "statuscategorychangedate": "2025-03-28T16:17:01.193-0700",
    "customfield_10350": null,
    "customfield_10352": null,
    "fixVersions": [],
    "customfield_10353": null,
    "customfield_11564": null,
    "customfield_10354": null,
    "customfield_11563": null,
    "customfield_11566": null,
    "customfield_10355": null,
    "customfield_10356": null,
    "customfield_11565": null,
    "customfield_10346": null,
    "customfield_10104": "2|hzv5hb:",
    "customfield_10347": null,
    "customfield_10348": null,
    "customfield_10349": null,
    "customfield_10340": null,
    "customfield_10341": null,
    "customfield_10100": "2025-04-15T11:06:19.762-0700",
    "customfield_10342": null,
    "priority": {
      "self": "https://thetalake.atlassian.net/rest/api/3/priority/2",
      "iconUrl": "https://thetalake.atlassian.net/images/icons/priorities/high_new.svg",
      "name": "High",
      "id": "2"
    },
    "customfield_10101": "1_*:*_2_*:*_2435222086_*|*_3_*:*_2_*:*_699425488_*|*_5_*:*_1_*:*_48545_*|*_10001_*:*_1_*:*_87826693",
    "customfield_10343": [],
    "customfield_10102": null,
    "customfield_10344": null,
    "customfield_10103": [
      {
        "id": 824,
        "name": "Zombie - Sprint 59",
        "state": "closed",
        "boardId": 10,
        "goal": "",
        "startDate": "2025-02-17T16:56:45.400Z",
        "endDate": "2025-03-26T07:46:25.000Z",
        "completeDate": "2025-03-31T08:52:03.437Z"
      },
      {
        "id": 956,
        "name": "Arrival - Sprint 60",
        "state": "closed",
        "boardId": 10,
        "goal": "",
        "startDate": "2025-03-31T08:52:43.337Z",
        "endDate": "2025-05-08T21:44:27.000Z",
        "completeDate": "2025-05-12T09:57:11.726Z"
      }
    ],
    "customfield_10345": null,
    "customfield_10335": null,
    "customfield_10336": null,
    "customfield_10337": null,
    "customfield_10338": null,
    "customfield_10339": null,
    "timeestimate": null,
    "versions": [],
    "status": {
      "self": "https://thetalake.atlassian.net/rest/api/3/status/6",
      "description": "The issue is considered finished, the resolution is correct. Issues which are closed can be reopened.",
      "iconUrl": "https://thetalake.atlassian.net/images/icons/statuses/closed.png",
      "name": "Issue Closed",
      "id": "6",
      "statusCategory": {
        "self": "https://thetalake.atlassian.net/rest/api/3/statuscategory/3",
        "id": 3,
        "key": "done",
        "colorName": "green",
        "name": "Done"
      }
    },
    "customfield_10330": null,
    "customfield_10452": null,
    "customfield_10331": null,
    "customfield_11663": null,
    "customfield_11662": null,
    "customfield_10453": null,
    "customfield_10332": null,
    "customfield_10333": null,
    "customfield_11664": null,
    "customfield_10334": null,
    "customfield_10324": null,
    "customfield_10445": null,
    "customfield_10325": null,
    "customfield_10446": null,
    "customfield_10326": null,
    "customfield_10447": null,
    "customfield_10448": null,
    "customfield_10327": null,
    "customfield_10328": null,
    "customfield_10449": null,
    "aggregatetimeestimate": null,
    "customfield_10329": null,
    "creator": {
      "self": "https://thetalake.atlassian.net/rest/api/3/user?accountId=6080b9dcb8927500729602d7",
      "accountId": "6080b9dcb8927500729602d7",
      "emailAddress": "kannan@thetalake.com",
      "avatarUrls": {
        "48x48": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "24x24": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "16x16": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png",
        "32x32": "https://secure.gravatar.com/avatar/2712d8f73af4dfe62dbcf561a1a5a7c7?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FKA-3.png"
      },
      "displayName": "Kannan Appachi",
      "active": true,
      "timeZone": "America/Los_Angeles",
      "accountType": "atlassian"
    },
    "aggregateprogress": {
      "progress": 0,
      "total": 0
    },
    "customfield_10320": null,
    "customfield_10321": null,
    "customfield_10443": null,
    "customfield_10322": null,
    "customfield_10323": null,
    "customfield_10444": null,
    "customfield_10434": null,
    "customfield_10313": null,
    "customfield_10314": null,
    "customfield_10435": null,
    "customfield_10315": null,
    "customfield_10436": null,
    "customfield_10437": null,
    "customfield_10316": null,
    "customfield_10438": null,
    "customfield_10439": null,
    "customfield_10319": null,
    "timespent": null,
    "aggregatetimespent": null,
    "customfield_10673": null,
    "customfield_10310": null,
    "customfield_10311": null,
    "customfield_10432": null,
    "customfield_10302": null,
    "customfield_10303": null,
    "customfield_10304": [],
    "customfield_10305": null,
    "customfield_10427": null,
    "customfield_10306": null,
    "customfield_10307": null,
    "customfield_10428": null,
    "customfield_10308": null,
    "customfield_10309": null,
    "workratio": -1,
    "created": "2025-02-17T09:05:31.290-0800",
    "customfield_10540": null,
    "customfield_10301": null,
    "customfield_10771": null,
    "customfield_10772": null,
    "customfield_10773": null,
    "customfield_10522": null,
    "customfield_10523": null,
    "security": null,
    "customfield_10524": null,
    "customfield_10525": null,
    "attachment": [],
    "customfield_11299": null,
    "customfield_10511": null,
    "customfield_10514": null,
    "customfield_10517": null,
    "customfield_10518": null,
    "customfield_10639": null
  }
}

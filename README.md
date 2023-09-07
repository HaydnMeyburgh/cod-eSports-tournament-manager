
# Call of Duty eSports Tournament Management system

This project is designed to help tournament organisers efficiently manage eSports events. It provides a range of features for creating, updating, and monitoring tournaments, matches, teams, and match results. In addition, the project also utilises WebSockets to broadcast live tournament updates and creations to client frontends for spectators and players.



## Table of Contents
- [Tech Stack](#tech-stack)
- [installation](#features)
- [API Reference](#api-reference)
- [For The Future](#for-the-future)


## Tech Stack

**Server:** Go, Gin

**Database:** Mongodb

**Communication:** WebSockets


## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/cod-esports-tournament-manager-backend.git
   cd cod-esports-tournament-manager-backend
   ```
2. Set up environment variables:

    ```bash
    SECRET_KEY=your_secret_key_here
    DATABASE_URL=your_database_url_here
    ```

3. Install Dependencies
    ```bash
    go mod tidy
    ```

4. Build and run the application
    ```bash
    go run main.go
    ```
    
## API Reference
### Users
#### Register a User

```http
  POST /users/register
```
**Request Body**
| Field | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required**. User's username |
| `email` | `string` | **Required**. User's Email |
| `password` | `string` | **Required**. User's password |


#### Login User

```http
  POST /users/login
```
**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | **Required**. User's Email |
| `password`      | `string` | **Required**. User's Password |

#### Logout User

```http
  POST /users/logout
```
**Security**: Cookie Token Authentication

#### Update User Profile
```http
  PUT /users/update/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the user to update |

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Optional**. New username |
| `Password`      | `string` | **Optional**. New password |

### Tournaments
#### Get All Tournaments
```http
  Get /tournaments
```
#### Get Tournament by ID
```http
  Get /tournaments/:id
```
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the tournament to fetch |

#### Create Tournament
```http
  POST /tournaments
```
**Security**: Cookie Token Authentication

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | **Required**. Tournament name |
| `description`      | `string` | **Optional**. tournament description |
| `start_date`      | `string` | **Required**. start date |
| `end_date`      | `string` | **Required**. end date |
| `organiser_id`      | `string` | **required**. organiser's id |
| `teams`      | `string` | **Optional**. teams participating |
| `matches`      | `string` | **Optional**. tournament matches |

#### Update Tournament
```http
  PUT /tournaments/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the tournament to update |

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | **Optional**. NewTournament name |
| `description`      | `string` | **Optional**. New tournament description |
| `start_date`      | `string` | **Optional**. New start date |
| `end_date`      | `string` | **Optional**. New end date |
| `teams`      | `string` | **Optional**. new teams participating |
| `matches`      | `string` | **Optional**. new tournament matches |

#### Delete Tournament
```http
  PUT /tournaments/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the tournament to delete |

### Matches
#### Get All Matches
```http
  Get /matches
```
#### Get Match by ID
```http
  Get /matches/:id
```
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match to fetch |

#### Create Match
```http
  POST /matches
```
**Security**: Cookie Token Authentication

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `tournament_id`      | `string` | **Optional**. Tournament id |
| `organiser_id`      | `string` | **Required**. organiser id |
| `team1_id`      | `string` | **Required**. team 1's id |
| `team2_id`      | `string` | **Required**. team 2's id |
| `date`      | `string` | **Optional**. date of match |
| `team1_name`      | `string` | **Optional**. team1's name |
| `team2_name`      | `string` | **Optional**. team2's name |

#### Update Match
```http
  PUT /matches/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match to update |

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `date`      | `string` | **Optional**. new date of match |
| `team1_name`      | `string` | **Optional**. new team1's name |
| `team2_name`      | `string` | **Optional**. new team2's name |

#### Delete Match
```http
  PUT /matches/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match to delete |

### Teams
#### Get All Teams
```http
  Get /teams
```
#### Get Team by ID
```http
  Get /teams/:id
```
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the team to fetch |

#### Create Team
```http
  POST /teams
```
**Security**: Cookie Token Authentication

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | **Required**. Team name |
| `organiser_id`      | `string` | **Required**. organiser id |
| `players`      | `string` | **Optional**. comma-separated list of player usernames |
| `tournament_id`      | `string` | **Optional**. team 2's id |

#### Update Team
```http
  PUT /teams/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the team to update |

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | **Optional**. New Team name |
| `players`      | `string` | **Optional**. New comma-separated list of player usernames |

#### Delete Team
```http
  PUT /teams/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the team to delete |

### Match Results
#### Get All Match Results
```http
  Get /match-results
```
#### Get Matchc Result by ID
```http
  Get /match-results/:id
```
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match result to fetch |

#### Create match Result
```http
  POST /match-results
```
**Security**: Cookie Token Authentication

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `match_id`      | `string` | **Required**. match id for which the result is related |
| `organiser_id`      | `string` | **Required**. organiser id |
| `winner_id`      | `string` | **Optional**. winning team's id |
| `loser_id`      | `string` | **Optional**. losing team's id |
| `winner_score`      | `int` | **Optional**. winning team's score |
| `loser_score`      | `int` | **Optional**. losing team's score |

#### Update Match Results
```http
  PUT /match-results/:id
```
**Security**: Cookie Token Authentication

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match result to update |

**Request Body**
| Field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `winner_score`      | `int` | **Optional**. New winning team's score |
| `loser_score`      | `int` | **Optional**. losing team's score |

#### Delete Match Result
```http
  PUT /match-results/:id
```
**Security**: Cookie Token Authentication
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of the match result to delete |

## For The Future
There are a couple things I would still like to add - I would like to add player profiles that contain stats and information about players that can be communicated through websocket to a client. I would like to create a feature that will create groups for teams, and then have the application auto-generate games based on the number of teams, and rules of the tournament (play every team once for example). It would also be nice to create bracket functionality that takes group standings and generates a playoff bracket.

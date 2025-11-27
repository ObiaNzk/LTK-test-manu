# LTK-test-manu

## Description 


## Endpoints

### POST - Events
Stores a new event on postgres

**Input**

    type body struct {
        Title       string    `json:"title" validate:"required"`
        Description string    `json:"description"`
        StartTime   time.Time `json:"start_time" validate:"required"`
        EndTime     time.Time `json:"end_time" validate:"required"`
    }

**Output**
    {
        "title": "pepito",
        "description": "pepito",
        "start_time": "12-09-18"
        "end_time": "12-12-2025"
        "created_at": "11-27-2025"
    }

**Status Code**
* 201 - Created
* 500 - Internal service error


### GET - Events
Returns all events in order (start_time)

Response
``` json [
    {
    "ID": "evt_001",
    "Title": "Reunión de equipo",
    "Description": "Revisión semanal de progreso",
    "StartTime": "2025-02-15T10:00:00Z",
    "EndTime": "2025-02-15T11:00:00Z",
    "CreatedAt": "2025-02-10T08:30:00Z"
    },
    {
    "ID": "evt_002",
    "Title": "Presentación de proyecto",
    "Description": "Presentación final al cliente",
    "StartTime": "2025-02-20T14:00:00Z",
    "EndTime": "2025-02-20T15:30:00Z",
    "CreatedAt": "2025-02-12T09:45:00Z"
    }
] ```

#!/bin/bash

docker exec -it nemoris-postgres psql -U nemoris -d nemoris -c "DELETE FROM reminders;"
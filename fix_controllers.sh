#!/bin/bash

# This script fixes the database collection initialization in all controller files

CONTROLLERS_DIR="/home/shams/D-Drive/Golang/Food-marketPlace/src/controllers"

# List of all controller files
for file in "$CONTROLLERS_DIR"/*.go; do
    if [[ -f "$file" ]]; then
        echo "Processing $file..."
        
        # Replace patterns like: var someCollection *mongo.Collection = database.OpenCollection(database.Client, "collectionname")
        # with lazy initialization
        
        sed -i 's/var \([a-zA-Z]*Collection\) \*mongo\.Collection = database\.OpenCollection(database\.Client, "\([^"]*\)")/var \1 *mongo.Collection\n\n\/\/ get\1 returns the \2 collection with lazy initialization\nfunc get\1() *mongo.Collection {\n\tif \1 == nil {\n\t\tclient := database.DBInstance()\n\t\t\1 = database.OpenCollection(client, "\2")\n\t}\n\treturn \1\n}/g' "$file"
        
    fi
done

echo "All controllers processed!"
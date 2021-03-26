HIER API ENDPOINTS DOKUMENTIEREN!

To Do:
 - Content Handler allgemein machen:
  - model 'ContentType', um eine Collection zu machen, wo dann wiederum Modelle für verschiedene Contenttypen enthalten sind:
   - Felder: ID, Type, Collection, Status (sowas wie veröffentlicht/unveröffentlicht) []Fields (Beim Ausgeben dann ein eigenes output struct, bei dem zumindest Collection nicht drin ist (wird nicht benötigt, da bei GET Frage mithilfe des Types in der 'ContentType' Collection nachgeschaut werden ka)nn, welche Collection dafür vorgesehen ist: Felder können auch geschachtelt werden) Initialisierung für "Blogposts" und "Events" (Funktion newContentType(typeName, collection, []fields) benutzen?)
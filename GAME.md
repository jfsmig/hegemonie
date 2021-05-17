# Hegemonie: The gameplay

## Civilization-like

The Hegemonie game engine is a technical basis for online RPG games. The
description of the game will sound familiar to RPG players:

* Each **user** manages characters
  
* Each **character** manages cities
  
* Each **city** produces **resources**, grows a **skills** tree, evolves
  with **building**, trains **troops**, controls **armies**, holds
  **artifacts**.
  
* Each **skill** brings some modifiers to the behavior of the city owning
  it, like altering of the production of resources on the city, allowing or
  forbidding other skills.
  
* Each **building** might also have an impact on the resource storage and
  production, but it can also
  
* **Armies** can be setup to gather troops, resources and artifacts. An army can
  execute a sequence of commands to move across the map, attack other cities,
  etc.
  
* **Artifacts** can be stored or hidden in cities, and transported by armies. As
  a consequence, they can be stolen and dropped by the armies. The primary goal
  of an artifact is to trigger quests evolutions upon the artifact lifetime.

There are plans to implement a quest system, destined to be triggered by NPC and
artifacts. But no code yet.

## Actually a 3X game, not a 4X :)

The game engine doesn't help building [4X](https://en.wikipedia.org/wiki/4X)
game instances. Particularly because there is no **exploration** in Hegemonie.
Each map is public and thus well-known. However, the experience on Hegemonie
proved that extermination, expansion and exploitation took a considerable part
of the gameplay for a large portion of users.

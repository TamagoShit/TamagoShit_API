-- PostgreSQL dump for TamagoShit database with all tables
-- Updated to include missing 'tamas' table

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS "Sponsor" CASCADE;
DROP TABLE IF EXISTS "Friends" CASCADE;
DROP TABLE IF EXISTS "tamas" CASCADE;
DROP TABLE IF EXISTS "Sickness" CASCADE;
DROP TABLE IF EXISTS "Trait" CASCADE;
DROP TABLE IF EXISTS "LifeChoices" CASCADE;
DROP TABLE IF EXISTS "Race" CASCADE;
DROP TABLE IF EXISTS "Malus" CASCADE;
DROP TABLE IF EXISTS "Event" CASCADE;
DROP TABLE IF EXISTS "Bonus" CASCADE;
DROP TABLE IF EXISTS "Tama_stats" CASCADE;
DROP TABLE IF EXISTS "Users" CASCADE;

-- Create sequences for auto-increment IDs
CREATE SEQUENCE bonus_bonusid_seq START 1;
CREATE SEQUENCE event_eventid_seq START 1;
CREATE SEQUENCE lifechoices_lifechoicesid_seq START 1;
CREATE SEQUENCE malus_malusid_seq START 1;
CREATE SEQUENCE race_raceid_seq START 1;
CREATE SEQUENCE sickness_sicknessid_seq START 1;
CREATE SEQUENCE tama_stats_tamastatid_seq START 1;
CREATE SEQUENCE tamas_tamaid_seq START 1;
CREATE SEQUENCE trait_traitid_seq START 1;
CREATE SEQUENCE users_userid_seq START 1;

-- Table: Bonus
CREATE TABLE "Bonus" (
  "BonusId" INTEGER NOT NULL DEFAULT nextval('bonus_bonusid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "Effect" VARCHAR(255),
  PRIMARY KEY ("BonusId")
);

CREATE INDEX "idx_bonus_name" ON "Bonus" ("Name");

-- Table: Event
CREATE TABLE "Event" (
  "EventId" INTEGER NOT NULL DEFAULT nextval('event_eventid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "Bonus" VARCHAR(255),
  "Malus" VARCHAR(255),
  PRIMARY KEY ("EventId")
);

CREATE INDEX "idx_event_name" ON "Event" ("Name");

-- Table: Users
CREATE TABLE "Users" (
  "UserId" INTEGER NOT NULL DEFAULT nextval('users_userid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "LastName" VARCHAR(100) NOT NULL,
  "UserName" VARCHAR(100) NOT NULL UNIQUE,
  "Email" VARCHAR(255) NOT NULL UNIQUE,
  "ProfilPicture" VARCHAR(255),
  "GamingTime" INTEGER DEFAULT 0,
  "CreationDate" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "LastConnectionDate" TIMESTAMP,
  "Password" VARCHAR(255) NOT NULL,
  "Role" VARCHAR(50) DEFAULT 'general',
  PRIMARY KEY ("UserId")
);

CREATE INDEX "idx_users_username" ON "Users" ("UserName");
CREATE INDEX "idx_users_email" ON "Users" ("Email");

-- Table: Tama_stats
CREATE TABLE "Tama_stats" (
  "TamaStatId" INTEGER NOT NULL DEFAULT nextval('tama_stats_tamastatid_seq'),
  "Food" INTEGER DEFAULT 0,
  "Play" INTEGER DEFAULT 0,
  "Cleaned" INTEGER DEFAULT 0,
  "CarAccident" INTEGER DEFAULT 0,
  "WorkAccident" INTEGER DEFAULT 0,
  "SocialSatis" DOUBLE PRECISION DEFAULT 0,
  "WorkSatis" DOUBLE PRECISION DEFAULT 0,
  "PersonalSatis" DOUBLE PRECISION DEFAULT 0,
  PRIMARY KEY ("TamaStatId")
);

-- Table: tamas
CREATE TABLE "tamas" (
  "TamaId" INTEGER NOT NULL DEFAULT nextval('tamas_tamaid_seq'),
  "UserId" INTEGER NOT NULL,
  "TamaStatsID" INTEGER NOT NULL,
  "Name" VARCHAR(100) NOT NULL,
  "Sexe" BOOLEAN,
  "Race" VARCHAR(100) NOT NULL,
  "Sickness" VARCHAR(255),
  "Birthday" TIMESTAMP,
  "DeathDay" TIMESTAMP,
  "Traits" VARCHAR(255),
  "LifeChoices" VARCHAR(255),
  PRIMARY KEY ("TamaId"),
  FOREIGN KEY ("UserId") REFERENCES "Users" ("UserId") ON DELETE CASCADE,
  FOREIGN KEY ("TamaStatsID") REFERENCES "Tama_stats" ("TamaStatId") ON DELETE CASCADE
);

CREATE INDEX "idx_tamas_userid" ON "tamas" ("UserId");
CREATE INDEX "idx_tamas_tamastatsid" ON "tamas" ("TamaStatsID");

-- Table: Friends
CREATE TABLE "Friends" (
  "UserID" INTEGER NOT NULL,
  "FriendID" INTEGER NOT NULL,
  "DateBecameFriends" DATE NOT NULL DEFAULT CURRENT_DATE,
  PRIMARY KEY ("UserID", "FriendID"),
  FOREIGN KEY ("UserID") REFERENCES "Users" ("UserId") ON DELETE CASCADE,
  FOREIGN KEY ("FriendID") REFERENCES "Users" ("UserId") ON DELETE CASCADE
);

CREATE INDEX "idx_friends_userid" ON "Friends" ("UserID");
CREATE INDEX "idx_friends_friendid" ON "Friends" ("FriendID");

-- Table: LifeChoices
CREATE TABLE "LifeChoices" (
  "LifeChoicesId" INTEGER NOT NULL DEFAULT nextval('lifechoices_lifechoicesid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "Traits" VARCHAR(255),
  PRIMARY KEY ("LifeChoicesId")
);

CREATE INDEX "idx_lifechoices_name" ON "LifeChoices" ("Name");

-- Table: Malus
CREATE TABLE "Malus" (
  "MalusId" INTEGER NOT NULL DEFAULT nextval('malus_malusid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "Effect" VARCHAR(255),
  PRIMARY KEY ("MalusId")
);

CREATE INDEX "idx_malus_name" ON "Malus" ("Name");

-- Table: Race
CREATE TABLE "Race" (
  "RaceId" INTEGER NOT NULL DEFAULT nextval('race_raceid_seq'),
  "Name" VARCHAR(100) NOT NULL UNIQUE,
  "Description" TEXT,
  "Bonus" VARCHAR(255),
  "Malus" VARCHAR(255),
  PRIMARY KEY ("RaceId")
);

CREATE INDEX "idx_race_name" ON "Race" ("Name");

-- Table: Sickness
CREATE TABLE "Sickness" (
  "SicknessId" INTEGER NOT NULL DEFAULT nextval('sickness_sicknessid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "ExpirationDays" INTEGER,
  "Bonus" VARCHAR(255),
  "Malus" VARCHAR(255),
  PRIMARY KEY ("SicknessId")
);

CREATE INDEX "idx_sickness_name" ON "Sickness" ("Name");

-- Table: Sponsor
CREATE TABLE "Sponsor" (
  "SponsorId" INTEGER NOT NULL,
  "SponsoredId" INTEGER NOT NULL,
  "DateOfSponsor" DATE NOT NULL DEFAULT CURRENT_DATE,
  PRIMARY KEY ("SponsorId", "SponsoredId"),
  FOREIGN KEY ("SponsorId") REFERENCES "Users" ("UserId") ON DELETE CASCADE,
  FOREIGN KEY ("SponsoredId") REFERENCES "Users" ("UserId") ON DELETE CASCADE
);

CREATE INDEX "idx_sponsor_sponsorid" ON "Sponsor" ("SponsorId");
CREATE INDEX "idx_sponsor_sponsoredid" ON "Sponsor" ("SponsoredId");

-- Table: Trait
CREATE TABLE "Trait" (
  "TraitId" INTEGER NOT NULL DEFAULT nextval('trait_traitid_seq'),
  "Name" VARCHAR(100) NOT NULL,
  "Description" TEXT,
  "Bonus" VARCHAR(255),
  "Malus" VARCHAR(255),
  PRIMARY KEY ("TraitId")
);

CREATE INDEX "idx_trait_name" ON "Trait" ("Name");

-- End of PostgreSQL dump
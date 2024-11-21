DROP SCHEMA IF EXISTS stew_accounts CASCADE;
CREATE SCHEMA stew_accounts;

CREATE TABLE stew_accounts.accounts
(
    "uuid"          uuid        NOT NULL,
    "name"          VARCHAR(16) NOT NULL,
    "gems"          BIGINT      NOT NULL DEFAULT 0,
    "coins"         BIGINT      NOT NULL DEFAULT 0,
    "lastLogin"     TIMESTAMP            DEFAULT NULL,
    "totalPlayTime" INT         NOT NULL DEFAULT 0,
    PRIMARY KEY ("uuid")
);

CREATE TABLE stew_accounts.accountAmplifierThank
(
    "amplifierId" BIGSERIAL UNIQUE NOT NULL,
    "playerUUID"  uuid             NOT NULL,
    "time"        TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("amplifierId", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.customData
(
    "id"   BIGSERIAL NOT NULL,
    "name" TEXT      NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.accountCustomData
(
    "playerUUID"   uuid   NOT NULL,
    "customDataId" BIGINT NOT NULL,
    "data"         BIGINT NOT NULL,
    PRIMARY KEY ("playerUUID", "customDataId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("customDataId") REFERENCES stew_accounts.customData ("id")
);

CREATE TABLE stew_accounts.accountFavouriteNano
(
    "playerUUID" uuid     NOT NULL,
    "gameId"     SMALLINT NOT NULL,
    PRIMARY KEY ("playerUUID", "gameId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountFriend
(
    "uuidSource" uuid         NOT NULL,
    "uuidTarget" uuid         NOT NULL,
    "status"     VARCHAR(100) NOT NULL,
    "favourite"  BOOLEAN      NOT NULL DEFAULT FALSE,
    PRIMARY KEY ("uuidSource", "uuidTarget"),
    CHECK ("uuidSource" < "uuidTarget"),
    FOREIGN KEY ("uuidSource") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("uuidTarget") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountFriendData
(
    "playerUUID" uuid     NOT NULL,
    "status"     SMALLINT NOT NULL,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountIgnore
(
    "uuidIgnorer" uuid NOT NULL,
    "uuidIgnored" uuid NOT NULL,
    PRIMARY KEY ("uuidIgnorer", "uuidIgnored"),
    CHECK ("uuidIgnorer" < "uuidIgnored"),
    FOREIGN KEY ("uuidIgnorer") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("uuidIgnored") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.items
(
    "id"     BIGSERIAL    NOT NULL,
    "name"   varchar(100) NOT NULL,
    "rarity" INT          NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.accountInventory
(
    "id"         BIGSERIAL NOT NULL,
    "playerUUID" uuid      NOT NULL,
    "itemId"     BIGINT    NOT NULL,
    "count"      BIGINT    NOT NULL,
    PRIMARY KEY ("id", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("itemId") REFERENCES stew_accounts.items ("id")
);

CREATE TABLE stew_accounts.accountKits
(
    "playerUUID" uuid             NOT NULL,
    "kitId"      BIGSERIAL UNIQUE NOT NULL,
    "active"     BOOLEAN          NOT NULL DEFAULT TRUE,
    PRIMARY KEY ("playerUUID", "kitId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.kitProgression
(
    "playerUUID"    uuid   NOT NULL,
    "kitId"         BIGINT NOT NULL,
    "xp"            BIGINT NOT NULL DEFAULT 0,
    "level"         BIGINT NOT NULL DEFAULT 0,
    "upgrade_level" BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY ("playerUUID", "kitId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("kitId") REFERENCES stew_accounts.accountKits ("kitId")
);

CREATE TABLE stew_accounts.stats
(
    "id"   BIGSERIAL    NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.accountKitStats
(
    "playerUUID" uuid   NOT NULL,
    "kitId"      BIGINT NOT NULL,
    "statId"     BIGINT NOT NULL,
    "value"      BIGINT NOT NULL,
    PRIMARY KEY ("playerUUID", "kitId", "statId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("kitId") REFERENCES stew_accounts.accountKits ("kitId"),
    FOREIGN KEY ("statId") REFERENCES stew_accounts.stats ("id")
);

CREATE TABLE stew_accounts.accountLevelReward
(
    "playerUUID" uuid   NOT NULL,
    "level"      BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.polls
(
    "id"          BIGSERIAL NOT NULL,
    "enabled"     BOOLEAN            DEFAULT FALSE,
    "question"    TEXT      NOT NULL,
    "answerA"     TEXT      NOT NULL,
    "answerB"     TEXT               DEFAULT NULL,
    "answerC"     TEXT               DEFAULT NULL,
    "answerD"     TEXT               DEFAULT NULL,
    "coinReward"  BIGINT    NOT NULL DEFAULT 0,
    "displayType" INT       NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.accountPolls
(
    "playerUUID" uuid    NOT NULL,
    "pollId"     BIGINT  NOT NULL,
    "value"      BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY ("playerUUID", "pollId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("pollId") REFERENCES stew_accounts.polls ("id")
);

CREATE TABLE stew_accounts.accountPunishments
(
    "id"               BIGSERIAL      NOT NULL,
    "playerUUID"       uuid           NOT NULL,
    "category"         TEXT           NOT NULL,
    "sentence"         TEXT           NOT NULL,
    "reason"           TEXT           NOT NULL,
    "time"             TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "duration"         NUMERIC(16, 2) NOT NULL,
    "adminUUID"        uuid           NOT NULL,
    "severity"         SMALLINT       NOT NULL DEFAULT 1,
    "removed"          BOOLEAN        NOT NULL DEFAULT false,
    "reasonOfRemoval"  TEXT           NOT NULL DEFAULT '',
    "removerAdminUUID" uuid           NOT NULL,
    PRIMARY KEY ("id", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("adminUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("removerAdminUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountRanks
(
    "playerUUID"     uuid NOT NULL,
    "rankIdentifier" VARCHAR(10) DEFAULT 'PLAYER',
    "primaryGroup"   BOOLEAN     DEFAULT true,
    PRIMARY KEY ("rankIdentifier", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountStatsAllTime
(
    "playerUUID" uuid   NOT NULL,
    "statId"     BIGINT NOT NULL,
    "value"      BIGINT NOT NULL,
    PRIMARY KEY ("statId", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("statId") REFERENCES stew_accounts.stats ("id")
);

CREATE TABLE stew_accounts.tasks
(
    "id"   BIGSERIAL    NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.accountTasks
(
    "id"         BIGSERIAL NOT NULL,
    "playerUUID" uuid      NOT NULL,
    "taskId"     BIGINT    NOT NULL,
    PRIMARY KEY ("id", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("taskId") REFERENCES stew_accounts.tasks ("id")
);

CREATE TABLE stew_accounts.accountThankTransactions
(
    "id"             BIGSERIAL NOT NULL,
    "receiverUUID"   uuid      NOT NULL,
    "senderUUID"     uuid      NOT NULL,
    "thankAmount"    BIGINT    NOT NULL DEFAULT 0,
    "reason"         TEXT      NOT NULL DEFAULT '',
    "ignoreCooldown" BOOLEAN   NOT NULL DEFAULT FALSE,
    "claimed"        BOOLEAN   NOT NULL DEFAULT FALSE,
    "sentTime"       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "claimTime"      TIMESTAMP,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("receiverUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("senderUUID") REFERENCES stew_accounts.accounts ("uuid")
);


CREATE OR REPLACE FUNCTION stew_accounts.addThank(
    IN inReceiverUUID uuid,
    IN inSenderUUID uuid,
    IN inThankAmount BIGINT,
    IN inReason TEXT,
    IN inIgnoreCooldown BOOLEAN,
    OUT success BOOLEAN)
AS
$$
DECLARE
    insertSuccess BOOLEAN := false;
    p_rows        BIGINT  := 0;
BEGIN
    INSERT INTO stew_accounts.accountThankTransactions
    ("receiverUUID", "senderUUID", "thankAmount", "reason", "ignoreCooldown")
    VALUES (inReceiverUUID, inSenderUUID, inThankAmount, inReason, inIgnoreCooldown);

    GET DIAGNOSTICS p_rows := ROW_COUNT;
    IF p_rows > 0 THEN
        insertSuccess := true;
    END IF;

    success := insertSuccess;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_accounts.checkAmplifierThank(IN inPlayerUUID uuid, IN inAmplifierId BIGINT, OUT canThank BOOLEAN)
AS
$$
DECLARE
    countValue INT;
BEGIN
    SELECT COUNT(*)
    INTO countValue
    FROM stew_accounts.accountAmplifierThank
    WHERE accountAmplifierThank."playerUUID" = inPlayerUUID
      AND accountAmplifierThank."amplifierId" = inAmplifierId;

    IF countValue > 0 THEN
        canThank := false;
    ELSE
        canThank := true;
        INSERT INTO stew_accounts.accountAmplifierThank ("playerUUID", "amplifierId")
        VALUES (inPlayerUUID, inAmplifierId);
    END IF;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_accounts.claimThank(IN inPlayerUUID uuid, OUT amountClaimed INT, OUT uniqueThank INT)
AS
$$
BEGIN

    SELECT SUM("thankAmount")
    INTO amountClaimed
    FROM stew_accounts.accountThankTransactions
    WHERE accountThankTransactions."receiverUUID" = inPlayerUUID
      AND claimed = false;

    UPDATE stew_accounts.accountThankTransactions
    SET claimed     = true,
        "claimTime" = CURRENT_TIMESTAMP
    WHERE accountThankTransactions."receiverUUID" = inPlayerUUID
      AND accountThankTransactions.claimed = false;

    SELECT COUNT(DISTINCT "senderUUID")
    INTO uniqueThank
    FROM stew_accounts.accountThankTransactions
    WHERE accountThankTransactions."receiverUUID" = inPlayerUUID
      AND claimed = false;

END
$$ LANGUAGE plpgsql;


CREATE TABLE stew_accounts.accountTitle
(
    "playerUUID" uuid         NOT NULL,
    "trackName"  VARCHAR(255) NOT NULL,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.accountWinstreak
(
    "playerUUID" uuid   NOT NULL,
    "gameId"     BIGINT NOT NULL,
    "value"      BIGINT NOT NULL,
    PRIMARY KEY ("playerUUID", "gameId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.botSpam
(
    "id"          BIGSERIAL NOT NULL,
    "text"        TEXT      NOT NULL,
    "punishments" BIGINT    NOT NULL,
    "enabled"     BOOLEAN   NOT NULL,
    "createdBy"   uuid      NOT NULL,
    "enabledBy"   uuid      NOT NULL,
    "disabledBy"  uuid      NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("createdBy") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("enabledBy") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("disabledBy") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.eloRating
(
    "playerUUID" uuid   NOT NULL,
    "gameId"     BIGINT NOT NULL,
    "elo"        BIGINT NOT NULL,
    PRIMARY KEY ("playerUUID", "gameId"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.incognitoStaff
(
    "playerUUID" uuid    NOT NULL,
    "status"     BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.npc
(
    "id"            BIGSERIAL      NOT NULL,
    "entityType"    TEXT           NOT NULL,
    "name"          TEXT           NOT NULL,
    "info"          TEXT     DEFAULT NULL,
    "world"         TEXT           NOT NULL,
    "x"             NUMERIC(16, 2) NOT NULL,
    "y"             NUMERIC(16, 2) NOT NULL,
    "z"             NUMERIC(16, 2) NOT NULL,
    "yaw"           SMALLINT DEFAULT 0,
    "pitch"         SMALLINT DEFAULT 0,
    "inHand"        TEXT           NOT NULL,
    "inHandData"    SMALLINT DEFAULT NULL,
    "helmet"        TEXT     DEFAULT NULL,
    "chestplate"    TEXT     DEFAULT NULL,
    "leggings"      TEXT     DEFAULT NULL,
    "boots"         TEXT     DEFAULT NULL,
    "metadata"      TEXT           NOT NULL,
    "skinValue"     TEXT     DEFAULT NULL,
    "skinSignature" TEXT     DEFAULT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.playerDisguiseName
(
    "playerUUID"  uuid        NOT NULL,
    "displayName" VARCHAR(16) NOT NULL,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.preferences
(
    "playerUUID" uuid    NOT NULL,
    "id"         BIGINT  NOT NULL,
    "value"      BOOLEAN NOT NULL,
    PRIMARY KEY ("playerUUID", "id"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.snapshots
(
    "id"          BIGSERIAL NOT NULL,
    "creatorUUID" uuid      NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("creatorUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.reportCategoryTypes
(
    "id"   SMALLINT    NOT NULL,
    "name" VARCHAR(16) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.reportTeams
(
    "id"   SMALLINT    NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.reportTeamMemberships
(
    "playerUUID" uuid     NOT NULL,
    "teamId"     SMALLINT NOT NULL,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("teamId") REFERENCES stew_accounts.reportTeams ("id")
);

CREATE TABLE stew_accounts.reports
(
    "id"           BIGSERIAL NOT NULL,
    "suspectUUID"  uuid      NOT NULL,
    "categoryId"   SMALLINT  NOT NULL,
    "snapshotId"   BIGSERIAL NOT NULL,
    "assignedTeam" BIGINT    NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("suspectUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("categoryId") REFERENCES stew_accounts.reportCategoryTypes ("id"),
    FOREIGN KEY ("snapshotId") REFERENCES stew_accounts.snapshots ("id"),
    FOREIGN KEY ("assignedTeam") REFERENCES stew_accounts.reportTeams ("id")
);

CREATE TABLE stew_accounts.reportHandlers
(
    "reportId"    BIGINT  NOT NULL,
    "handlerUUID" uuid    NOT NULL,
    "aborted"     BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY ("reportId"),
    FOREIGN KEY ("reportId") REFERENCES stew_accounts.reports ("id"),
    FOREIGN KEY ("handlerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.reportReasons
(
    "reportId"     BIGSERIAL   NOT NULL,
    "reporterUUID" uuid        NOT NULL,
    "reason"       TEXT        NOT NULL,
    "server"       VARCHAR(30) NOT NULL,
    "weight"       BIGINT      NOT NULL DEFAULT 0,
    "time"         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("reportId"),
    FOREIGN KEY ("reportId") REFERENCES stew_accounts.reports ("id"),
    FOREIGN KEY ("reporterUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.reportResultTypes
(
    "id"         SMALLINT    NOT NULL,
    "globalStat" BOOLEAN     NOT NULL,
    "name"       VARCHAR(16) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.reportResults
(
    "resultId"   BIGSERIAL NOT NULL,
    "reportId"   BIGINT    NOT NULL,
    "reason"     TEXT      NOT NULL,
    "closedTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("resultId"),
    FOREIGN KEY ("reportId") REFERENCES stew_accounts.reports ("id")
);

CREATE TABLE stew_accounts.snapshotTypes
(
    "id"   SMALLINT    NOT NULL,
    "name" VARCHAR(25) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_accounts.snapshotMessages
(
    "messageId"    BIGSERIAL   NOT NULL,
    "senderUUID"   uuid        NOT NULL,
    "server"       VARCHAR(30) NOT NULL,
    "time"         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "message"      TEXT        NOT NULL,
    "snapshotType" SMALLINT    NOT NULL,
    PRIMARY KEY ("messageId"),
    FOREIGN KEY ("senderUUID") REFERENCES stew_accounts.accounts ("uuid"),
    FOREIGN KEY ("snapshotType") REFERENCES stew_accounts.snapshotTypes ("id")
);

CREATE TABLE stew_accounts.snapshotRecipients
(
    "messageId"     BIGINT NOT NULL,
    "recipientUUID" uuid   NOT NULL,
    PRIMARY KEY ("messageId", "recipientUUID"),
    FOREIGN KEY ("recipientUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.snapshotMessageMap
(
    "snapshotId" BIGINT NOT NULL,
    "messageId"  BIGINT NOT NULL,
    PRIMARY KEY ("snapshotId", "messageId"),
    FOREIGN KEY ("snapshotId") REFERENCES stew_accounts.snapshots ("id"),
    FOREIGN KEY ("messageId") REFERENCES stew_accounts.snapshotMessages ("messageId")
);

CREATE TABLE stew_accounts.twofactor
(
    "playerUUID" uuid        NOT NULL,
    "secretKey"  VARCHAR(16) NOT NULL,
    PRIMARY KEY ("playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

CREATE TABLE stew_accounts.twofactorHistory
(
    "id"         BIGSERIAL   NOT NULL,
    "playerUUID" uuid        NOT NULL,
    "ipAddress"  VARCHAR(72) NOT NULL,
    "time"       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_accounts.accounts ("uuid")
);

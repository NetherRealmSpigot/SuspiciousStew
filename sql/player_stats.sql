DROP SCHEMA IF EXISTS stew_player_stats CASCADE;
CREATE SCHEMA IF NOT EXISTS stew_player_stats;

CREATE TABLE stew_player_stats.ipInfo
(
    "id"        BIGSERIAL   NOT NULL,
    "ipAddress" VARCHAR(72) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE stew_player_stats.playerInfo
(
    "uuid"    uuid        NOT NULL,
    "name"    VARCHAR(16) NOT NULL,
    "version" SMALLINT    NOT NULL,
    PRIMARY KEY ("uuid")
);

CREATE TABLE stew_player_stats.playerIps
(
    "playerUUID" uuid      NOT NULL,
    "ipInfoId"   BIGINT    NOT NULL,
    "date"       TIMESTAMP NOT NULL,
    FOREIGN KEY ("playerUUID") REFERENCES stew_player_stats.playerInfo ("uuid"),
    FOREIGN KEY ("ipInfoId") REFERENCES stew_player_stats.ipInfo ("id")
);

CREATE TABLE stew_player_stats.playerLoginSessions
(
    "id"         BIGSERIAL NOT NULL,
    "playerUUID" uuid      NOT NULL,
    "loginTime"  TIMESTAMP NOT NULL,
    "timeInGame" INT       NOT NULL DEFAULT 0,
    PRIMARY KEY ("id", "playerUUID"),
    FOREIGN KEY ("playerUUID") REFERENCES stew_player_stats.playerInfo ("uuid")
);

CREATE TABLE stew_player_stats.playerUniqueLogins
(
    "playerUUID" uuid      NOT NULL,
    "date"       TIMESTAMP NOT NULL,
    FOREIGN KEY ("playerUUID") REFERENCES stew_player_stats.playerInfo ("uuid")
);


CREATE OR REPLACE FUNCTION stew_player_stats.add_player_info(
    IN p_uuid uuid, IN p_name VARCHAR(16), IN p_version SMALLINT
) RETURNS VOID AS
$$
BEGIN
    INSERT INTO stew_player_stats.playerInfo (uuid, name, version) VALUES (p_uuid, p_name, p_version);
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.get_player_info(
    IN p_uuid uuid
) RETURNS SETOF stew_player_stats.playerInfo AS
$$
BEGIN
    RETURN QUERY SELECT * FROM stew_player_stats.playerInfo WHERE playerInfo.uuid = p_uuid;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.update_player_info(
    IN p_uuid uuid, IN p_name VARCHAR(16), IN p_version SMALLINT
) RETURNS VOID AS
$$
BEGIN
    UPDATE stew_player_stats.playerInfo SET name = p_name, version = p_version WHERE playerInfo.uuid = p_uuid;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.handle_player_logins(
    IN p_playerUUID uuid, IN p_ipInfoId BIGINT
) RETURNS VOID AS
$$
DECLARE
    currentTime TIMESTAMP := CURRENT_TIMESTAMP;
BEGIN
    INSERT INTO stew_player_stats.playerIps ("playerUUID", "ipInfoId", "date")
    VALUES (p_playerUUID, p_ipInfoId, currentTime)
    ON CONFLICT DO NOTHING;

    INSERT INTO stew_player_stats.playerUniqueLogins ("playerUUID", "date")
    VALUES (p_playerUUID, currentTime)
    ON CONFLICT DO NOTHING;

    INSERT INTO stew_player_stats.playerLoginSessions ("playerUUID", "loginTime")
    VALUES (p_playerUUID, currentTime)
    ON CONFLICT DO NOTHING;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.add_ip_info(
    IN p_ipAddress VARCHAR(72)
) RETURNS VOID AS
$$
BEGIN
    INSERT INTO stew_player_stats.ipInfo ("ipAddress") VALUES (p_ipAddress);
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.get_ip_info(
    IN p_ipAddress VARCHAR(72)
) RETURNS SETOF stew_player_stats.ipInfo AS
$$
BEGIN
    RETURN QUERY SELECT * FROM stew_player_stats.ipInfo WHERE ipInfo."ipAddress" = p_ipAddress;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.get_session_id(
    IN p_uuid uuid
) RETURNS SETOF stew_player_stats.playerLoginSessions AS
$$
BEGIN
    RETURN QUERY SELECT *
                 FROM stew_player_stats.playerLoginSessions
                 WHERE playerLoginSessions."playerUUID" = p_uuid;
END
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION stew_player_stats.update_login_session(
    IN p_id BIGINT
) RETURNS VOID AS
$$
BEGIN
    UPDATE stew_player_stats.playerLoginSessions
    SET "timeInGame" = EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - "loginTime")) / 60
    WHERE playerLoginSessions.id = p_id;
END
$$ LANGUAGE plpgsql;
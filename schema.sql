
CREATE TABLE "comment" (
    content  VARCHAR 

    ,
    "Date"   VARCHAR 

    ,
    id       VARCHAR 

    ,
    user_id  VARCHAR 

     NOT NULL,
    post_id  VARCHAR 

     NOT NULL
);

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE post (
    id          VARCHAR 

     NOT NULL,
    title       VARCHAR 

    ,
    content     VARCHAR 

    ,
    "Date"      VARCHAR 

    ,
    commenting  VARCHAR 

    ,
    user_id     VARCHAR 

     NOT NULL
);

ALTER TABLE post ADD CONSTRAINT post_pk PRIMARY KEY ( id );

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE "user" (
    id    VARCHAR 

     NOT NULL,
    name  VARCHAR 

                    
    );

ALTER TABLE "user" ADD CONSTRAINT user_pk PRIMARY KEY ( id );

ALTER TABLE "comment"
    ADD CONSTRAINT comment_post_fk FOREIGN KEY ( post_id )
        REFERENCES post ( id );

ALTER TABLE "comment"
    ADD CONSTRAINT comment_user_fk FOREIGN KEY ( user_id )
        REFERENCES "user" ( id );

ALTER TABLE post
    ADD CONSTRAINT post_user_fk FOREIGN KEY ( user_id )
        REFERENCES "user" ( id );



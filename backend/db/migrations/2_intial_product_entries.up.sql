INSERT INTO users (id, name, email, password) VALUES(1, "max", "max@cmail.com", "pw");
INSERT INTO users (name, email, password) VALUES("ann", "ann@cmail.com", "passwd"); 
INSERT INTO users (name, email, password) VALUES("bea", "bea@cmail.com", "pword");
INSERT INTO users (name, email, password) VALUES("tom", "tom@cmail.com", "password");

INSERT INTO products (id, name, description, price, duration) VALUES(1, "Tree Climbing Class - For Beginners", "Ever wanted to learn how to climb a tree but didn't know how to get started? Well good news, we have a course designed especially for you. That's right! In just 2 weeks you will be able to pat a koala in its natual habitat.", 10000, 14);
INSERT INTO products (id, name, description, price, duration) VALUES(2, "Get a 6-pack while drinking a 6-pack", "Do you often feel more tempted by a beer than a work out? Well we have designed a new work out that you can do while at the bar. We have fanshioned 6 drinking games that will also get you ripped.", 5000, 7);

INSERT INTO subscriptions (user_id, product_id, start_date, status) VALUES(1, 1, NOW(), "ACTIVE");
INSERT INTO subscriptions (user_id, product_id, start_date, status) VALUES(1, 2, NOW(), "ACTIVE");

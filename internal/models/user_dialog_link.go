package models

type UsersDialogLink struct {
	ID       int
	UserID   int
	DialogID int
	LinkName string
}

//CREATE TABLE users_dialogs_links (
//id SERIAL PRIMARY KEY,
//user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
//dialog_id INTEGER NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
//link_name TEXT,
//UNIQUE (user_id, dialog_id)
//);

//CREATE OR REPLACE FUNCTION set_default_link_name()
//RETURNS TRIGGER AS $$
//DECLARE
//other_user_login TEXT;
//BEGIN
//IF NEW.link_name IS NULL THEN
//-- Ищем login второго участника диалога (используем поле login!)
//SELECT u.login INTO other_user_login
//FROM users u
//JOIN user_dialogs_links udl ON u.id = udl.user_id
//WHERE udl.dialog_id = NEW.dialog_id
//AND udl.user_id != NEW.user_id
//LIMIT 1;
//
//NEW.link_name := other_user_login;
//END IF;
//
//RETURN NEW;
//END;
//$$ LANGUAGE plpgsql;
//
//CREATE TRIGGER trg_set_default_link_name
//BEFORE INSERT ON user_dialogs_links
//FOR EACH ROW EXECUTE FUNCTION set_default_link_name();

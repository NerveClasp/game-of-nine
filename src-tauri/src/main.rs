#![cfg_attr(
  all(not(debug_assertions), target_os = "windows"),
  windows_subsystem = "windows"
)]

use std::{sync::{Arc, Mutex}, time::Duration};

use tauri::{State, Manager};
// use tokio::time::sleep;

#[derive(Default)]
struct Counter(Arc<Mutex<i32>>);

// #[derive(Default)]
// struct Game {
//   started: bool,
//   players: Vec<Player>
// }
//
// #[derive(Default)]
// struct GameState(Arc<Mutex<Game>>);
//
// #[derive(Default)]
// struct Player {
//   name: String,
//   money: i32,
//   is_computer: bool
// }

fn main() {
  tauri::Builder::default()
    .manage(Counter(Default::default()))
    .setup(|app|{
      let app_handle = app.app_handle();
      // let mut game = GameState::default();

      tauri::async_runtime::spawn(async move {
        app_handle.listen_global("add-player", |event| {
          println!("got add-player with payload {:?}", event.payload())
        });

        // loop {
        //   sleep(Duration::from_millis(2000)).await;
        //   println!("sending backend ping");
        //   app_handle.emit_all("backend-ping", "ping").unwrap();
        // }
      });

        Ok(())
    })
    .invoke_handler(tauri::generate_handler![
      hello_world,
      counter_inc,
    ])
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}

#[tauri::command]
fn hello_world() -> String {
  "Hello world".to_string()
}

#[tauri::command]
fn counter_inc(num: i32, counter: State<'_, Counter>) -> String {
  let mut val = counter.0.lock().unwrap();
  *val += num;
  format!("{val}")
}

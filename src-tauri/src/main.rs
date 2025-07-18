#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::sync::{Arc, Mutex};

use tauri::{Manager, State};

#[derive(Default)]
struct Counter(Arc<Mutex<i32>>);

fn main() {
    tauri::Builder::default()
        .manage(Counter(Default::default()))
        .setup(|app| {
            let app_handle = app.app_handle();

            tauri::async_runtime::spawn(async move {
                app_handle.listen_global("add-player", |event| {
                    println!("got add-player with payload {:?}", event.payload())
                });
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![hello_world, counter_inc,])
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

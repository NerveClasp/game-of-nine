#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::{
    sync::{Arc, Mutex},
    time::Duration,
};

use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;
use tauri::{Manager, State};

fn handle_read(mut stream: &TcpStream) {
    let mut buf = [0u8; 4096];
    match stream.read(&mut buf) {
        Ok(_) => {
            let req_str = String::from_utf8_lossy(&buf);
            println!("{}", req_str);
        }
        Err(e) => println!("Unable to read stream: {}", e),
    }
}

fn handle_write(mut stream: TcpStream) {
    let response = b"HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n<html><body>Hello world</body></html>\r\n";
    match stream.write(response) {
        Ok(_) => println!("Response sent"),
        Err(e) => println!("Failed sending response: {}", e),
    }
}

fn handle_client(stream: TcpStream) {
    handle_read(&stream);
    handle_write(stream);
}
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
    let listener = TcpListener::bind("127.0.0.1:8765").unwrap();
    println!("Listening for connections on port {}", 8765);

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(|| handle_client(stream));
            }
            Err(e) => {
                println!("Unable to connect: {}", e);
            }
        }
    }
    tauri::Builder::default()
        .manage(Counter(Default::default()))
        .setup(|app| {
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

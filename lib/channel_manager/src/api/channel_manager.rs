use bioenpro4to_channel_manager::channels::root_channel::RootChannel;
use std::ptr::{null, null_mut};
use crate::api::utils::{ChannelInfo, RawPacket, KeyNonce};
use std::os::raw::c_char;
use std::ffi::{CStr, CString};
use tokio::runtime::Runtime;
use bioenpro4to_channel_manager::channels::actor_channel::DailyChannelManager;
use crate::api::int_to_category;

#[no_mangle]
pub extern "C" fn new_root_channel(mainnet: usize) -> *mut RootChannel{
    let mainnet = match mainnet{
        0 => false,
        _ => true
    };
    let root = RootChannel::new(mainnet);
    Box::into_raw(Box::new(root))
}

#[no_mangle]
pub unsafe extern "C" fn import_root_channel_from_tangle(channel_info: *const ChannelInfo, channel_psw: *const c_char, mainnet: usize) -> *mut RootChannel{
    let mainnet = match mainnet{
        0 => false,
        _ => true
    };

    let channel_info = match channel_info.as_ref(){
        None => return null_mut(),
        Some(info) => info
    };
    let channel_info = match channel_info.to_ch_info(){
        Ok(info) => info,
        Err(_) => return null_mut()
    };

    let channel_psw = match CStr::from_ptr(channel_psw).to_str(){
        Ok(psw) => psw,
        Err(_) => return null_mut()
    };

    Runtime::new().unwrap().block_on(async {
        match RootChannel::import_from_tangle(
            channel_info.channel_id(),
            channel_info.announce_id(),
            channel_psw,
            mainnet).await
        {
            Ok(root) => Box::into_raw(Box::new(root)),
            Err(_) => null_mut()
        }
    })
}

#[no_mangle]
pub unsafe extern "C" fn drop_root_channel(channel: *mut RootChannel){
    channel.drop_in_place();
}

#[no_mangle]
pub unsafe extern "C" fn open_root_channel(channel: *mut RootChannel, channel_psw: *const c_char) -> *const ChannelInfo{
    let ch = match channel.as_mut(){
        None => return null(),
        Some(ch) => ch
    };

    let state_psw = match CStr::from_ptr(channel_psw).to_str(){
        Ok(state_psw) => state_psw,
        Err(_) => return null()
    };

    Runtime::new().unwrap().block_on(async {
      match ch.open(state_psw).await{
          Ok(info) => ChannelInfo::from_ch_info(info),
          Err(_) => null()
      }
    })
}

#[no_mangle]
pub unsafe extern "C" fn get_create_daily_actor_channel(
    channel: *mut RootChannel, category: usize, actor_id: *const c_char,
    state_psw: *const c_char, day: u16, month: u16, year: u16
) -> *mut DailyChannelManager{
    let ch = match channel.as_mut(){
        None => return null_mut(),
        Some(ch) => ch
    };

    let category = match int_to_category(category){
        Ok(c) => c,
        Err(_) => return null_mut()
    };

    let actor_id = match CStr::from_ptr(actor_id).to_str(){
        Ok(state_psw) => state_psw,
        Err(_) => return null_mut()
    };

    let state_psw = match CStr::from_ptr(state_psw).to_str(){
        Ok(state_psw) => state_psw,
        Err(_) => return null_mut()
    };

    Runtime::new().unwrap().block_on(async {
        match ch.get_or_create_daily_actor_channel(category, actor_id, state_psw, day, month, year).await{
            Ok(ch) => Box::into_raw(Box::new(ch)),
            Err(_) => null_mut()
        }
    })
}

#[no_mangle]
pub unsafe extern "C" fn drop_daily_channel_manager(channel: *mut DailyChannelManager){
    channel.drop_in_place();
}

#[no_mangle]
pub unsafe extern "C" fn print_channel_tree(root: *mut RootChannel){
    let root = match root.as_mut() {
        None => { return; }
        Some(root) => root
    };
    root.print_nested_channel_info();
}

#[no_mangle]
pub unsafe extern "C" fn send_raw_packet(root: *mut DailyChannelManager, packet: *const RawPacket, key_nonce: *const KeyNonce) -> *const c_char{
    let root = root.as_mut();
    let p = packet.as_ref();
    let kn = key_nonce.as_ref();

    match (&root, &p){
        (None, _) => return null(),
        (_, None) => return null(),
        _ => {}
    };

    let root = root.unwrap();
    let p = p.unwrap();
    let public = p.public();
    let masked = p.masked();
    let opt_kn = match kn{
        None => None,
        Some(kn) => Some((kn.key.clone(), kn.nonce.clone()))
    };

    let res = Runtime::new().unwrap().block_on(async {
        root.send_raw_packet(public, masked, opt_kn).await
    });

    match res{
        Ok(res) => CString::new(res).map_or(null(), |h| h.into_raw()),
        Err(_) => null()
    }
}

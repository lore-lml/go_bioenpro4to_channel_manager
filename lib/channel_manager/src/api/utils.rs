use std::os::raw::c_char;
use bioenpro4to_channel_manager::channels::ChannelInfo as ChInfo;
use std::ffi::CString;

#[repr(C)]
pub struct ChannelInfo{
    pub channel_id: *const c_char,
    pub announce_id: *const c_char
}

impl ChannelInfo{
    pub fn from_ch_info(ch_info: ChInfo) -> Self{
        let channel_id = CString::new(ch_info.channel_id()).unwrap().into_raw();
        let announce_id = CString::new(ch_info.announce_id()).unwrap().into_raw();
        ChannelInfo{ channel_id, announce_id }
    }
}

#[no_mangle]
pub unsafe extern "C" fn drop_channel_info(info: *mut ChannelInfo){
    Box::from_raw(info);
}

import { ConfigMix, Blacklist, CcCore } from "./out"


const a = ConfigMix.parse(`
    
token 123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
proxy http://127.0.0.1:9999 -enable
record record.json
groups -100123456789 -100987654321
log_channel 987654321
main_site zh.wikipedia.org
oauth_auth_url https://telegram-auth-bot.toolforge.org/auth?id={telegram_id}
oauth_query_url https://telegram-auth-bot.toolforge.org/query
oauth_query_key AAAAAAAAAAAAAA
wiki_list zhwiki
blacklist
message-start "入群门槛：在任意一个维基媒体计划网站注册超过 7 日且编辑 50 次以上。<b>不要为了入群而用快速编辑积累编辑次数，您会因此遭到封禁而无法再编辑。</b>\\n\\n/confirm 验证维基媒体账户\\n/deconfirm 解除与维基媒体账户的关联\\n/policy 查看机器人说明" 
message-policy "入群门槛：在任意一个维基媒体计划网站注册超过 7 日且编辑 50 次以上。<b>不要为了入群而用快速编辑积累编辑次数，您会因此遭到封禁而无法再编辑。</b>\\n\\n若要开始验证，请发送 /confirm 并按提示操作。机器人借由OAuth确认您的身份，并会检查您是否达到入群门槛。验证账户后，您就可以在群组中发言。您可以随时解除与站内账号的关联，若如此做，则机器人也会禁止您在群里发言。\\n\\n机器人在成功验证或解除关联后，会在一个日志频道记录这些操作。在群组中，可以通过指令查看其他用户对应的维基媒体用户名。\\n\\n机器人会记录的信息为：您的 Telegram 账户 1） 是否完成验证，2）是否正在验证中，3）Telegram ID，4）对应的维基媒体账号，5）完成验证的时间，6）上一次被群管禁言的期限" 
message-insufficient_right "请授予我 Ban Users 权限以便正常运作，感谢🙏" 
message-general_prompt "使用方法：指令 ID 备注" 
message-telegram_id_error "这里只接受数字 ID。" 
message-restore_silence "已按<a href=\\"tg://user?id={tg_id}\\">此用户</a> (<code>{tg_id}</code>) 先前的禁言记录实施禁言，请复查。" 
message-confirm_already "您已成功验证站内账户 {wp_name}。若无法发言，请联络群组管理员。若要更改关联的维基百科账户，请先使用 /deconfirm 解除关联，然后重新验证。" 
message-confirm_other_tg "已有其他Telegram账户验证为站内账户 {wp_name}，若要更改关联的维基百科账户，请先解除该Telegram账户的验证后重新验证。若无法自行解除，请联络群组管理员" 
message-confirm_conflict "您提供的维基百科用户名已验证为其他 Telegram 账户。" 
message-confirm_checking "正在检查您的维基百科账户" 
message-confirm_user_not_found "未找到 id 为 {mw_id} 的维基百科用户。" 
message-confirm_button "确认" 
message-confirm_wait "请点<a href=\\"{link}\\">此链接</a>按提示完成验证。\\n\\n完成后点击确认按钮。" 
message-confirm_confirming "您目前正在验证中。若上一个验证已无法继续，可以按下确认按钮结束验证，然后重新使用 /confirm 指令开始验证。" 
message-confirm_ineligible "对不起，您尚未达到入群门槛。\\n\\n<b>不要为了入群而用快速编辑积累编辑次数，您会因此遭到封禁而无法再编辑。</b>\\n\\n" 
message-confirm_session_lost "对不起，为确保验证有效，请重新使用 /confirm 指令进行验证。" 
message-confirm_complete "验证成功。" 
message-confirm_failed "验证失败，请使用 /policy 查看验证通过的条件。您可以在日后重新使用 /confirm 指令进行验证。若您确信您已满足条件而无法验证通过，请联系群组管理员。" 
message-confirm_log "#新 #u_{tg_id}\\n<a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 验证为 <a href=\\"https://{site}/wiki/Special:Contributions/{wp_name}\\">{wp_name}</a>" 
message-deconfirm_prompt "您可以使用下方的按钮来解除 Telegram 账户与维基百科账户的关联。解除关联后，您将无法在群内发言。" 
message-deconfirm_button "解除关联" 
message-deconfirm_succ "已解除与维基百科账户的关联。" 
message-deconfirm_not_confirmed "您目前没有验证维基百科用户身份。" 
message-deconfirm_log "#解 #u_{tg_id}\\n<a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 已解除与 <a href=\\"https://{site}/wiki/Special:Contributions/{wp_name}\\">{wp_name}</a> 的关联" 
message-new_member_hint "<a href=\\"tg://user?id={tg_id}\\">{tg_name}</a> (<code>{tg_id}</code>) 您好，请私聊我验证您的维基百科账号以取得发言权限。" 
message-add_whitelist_prompt "使用方法：/add_whitelist 用户ID 备注" 
message-add_whitelist_succ "<code>{tg_id}</code> 已加入白名单。" 
message-add_whitelist_log "#白 #u_{tg_id}\\n{adder} 已将 <a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 加入白名单，备注：{reason}" 
message-remove_whitelist_prompt "使用方法：/remove_whitelist 用户ID" 
message-remove_whitelist_not_found "未在白名单中找到此人" 
message-remove_whitelist_log "#白 #u_{tg_id}\\n{remover} 已将 <a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 移出白名单" 
message-remove_whitelist_succ "<code>{tg_id}</code> 已移出白名单。" 
message-whois_head "{name} (<code>{tg_id}</code>)\\n" 
message-whois_prompt "使用方法：\\n1. 以 /whois 回复要查询的用户\\n2. /whois <Telegram 数字ID>\\n3. /whois <站内用户名>" 
message-whois_not_found "未查到该用户。" 
message-whois_self "这是我自己" 
message-whois_bot "这是机器人" 
message-whois_has_mw "维基百科账号：<a href=\\"https://{site}/wiki/Special:Contributions/{wp_id}\\">{wp_id}</a>（于 {ctime} (UTC) 验证）" 
message-whois_no_mw "未验证维基百科账户\\n" 
message-whois_whitelisted "该用户在白名单中，备注：{reason}" 
message-whois_tg_name_unavailable "（无法获取Telegram用户名）" 
message-refuse_log "#禁 #u_{tg_id}\\n{refuser} 已禁止 <a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 进行验证。" 
message-accept_log "#禁 #u_{tg_id}\\n{acceptor} 已允许 <a href=\\"tg://user?id={tg_id}\\">{tg_id}</a> 进行验证。" 
message-lift_restriction_alert "{name} (<code>{tg_id}</code>) 被允许发言" 
message-silence_alert "{name} (<code>{tg_id}</code>) 被禁止发言" 
message-enable "启用成功" 
message-disable "禁用成功" 
message-enable_log "#开 #u_{tg_id}\\n群管 {enabler} 在 <a href=\\"{chat_link}\\">{chat_name}</a> 启用验证" 
message-disable_log "#关 #u_{tg_id}\\n群管 {enabler} 在 <a href=\\"{chat_link}\\">{chat_name}</a> 禁用验证"
    
`);

console.log(a);



// const b = ConfigMix.parse(a);
// console.log(b);

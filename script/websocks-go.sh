##################################################
# Anything wrong? Find me via telegram: @CN_SZTL #
##################################################

#! /bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

function set_fonts_colors(){
# Font colors
default_fontcolor="\033[0m"
red_fontcolor="\033[31m"
green_fontcolor="\033[32m"
# Background colors
green_backgroundcolor="\033[42;37m"
# Fonts
error_font="${red_fontcolor}[Error]${default_fontcolor}"
ok_font="${green_fontcolor}[OK]${default_fontcolor}"
}

function check_os(){
	clear
	echo -e "正在检测当前是否为ROOT用户..."
	if [[ $EUID -ne 0 ]]; then
		clear
		echo -e "${error_font}当前并非ROOT用户，请先切换到ROOT用户后再使用本脚本。"
		exit 1
	else
		clear
		echo -e "${ok_font}检测到当前为Root用户。"
	fi
	clear
	echo -e "正在检测此OS是否被支持..."
	if [ ! -z "$(cat /etc/issue | grep Debian)" ];then
		OS='debian'
		clear
		echo -e "${ok_font}该脚本支持您的系统。"
	elif [ ! -z "$(cat /etc/issue | grep Ubuntu)" ];then
		OS='ubuntu'
		clear
		echo -e "${ok_font}该脚本支持您的系统。"
	else
		clear
		echo -e "${error_font}目前暂不支持您使用的操作系统，请切换至Debian/Ubuntu。"
		exit 1
	fi
	clear
	echo -e "正在检测系统架构是否被支持..."
	system_bit=$(uname -m)
	if  [[ ${system_bit} = "x86_64" ]]; then
		clear
		echo -e "${ok_font}该脚本支持您的系统架构。"
	elif [[ ${system_bit} = "i386" ]]; then
		clear
		echo -e "${ok_font}该脚本支持您的系统架构。"
	elif [[ ${system_bit} = "i686" ]]; then
		clear
		echo -e "${ok_font}该脚本支持您的系统架构。"
	elif [[ ${system_bit} = "x86" ]]; then
		clear
		echo -e "${ok_font}该脚本支持您的系统架构。"
	elif [[ ${system_bit} = "arm64" ]]; then
		clear
		echo -e "${ok_font}该脚本支持您的系统架构。"
	else
		clear
		echo -e "${error_font}目前暂不支持您使用的系统架构，推荐切换至Debian/Ubuntu x86_64。"
		exit 1
	fi
}

function check_install_status(){
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [[ ${install_type} = "" ]]; then
		install_status="${red_fontcolor}未安装${default_fontcolor}"
		websocks_start_command="${red_fontcolor}未安装${default_fontcolor}"
	else
		install_status="${green_fontcolor}已安装${default_fontcolor}"
		websocks_start_command="${green_backgroundcolor}$(cat /usr/local/websocks/run_command.txt)${default_fontcolor}"
	fi
	websocks_program=$(find /usr/local/websocks/websocks)
	if [[ ${websocks_program} = "" ]]; then
		websocks_status="${red_fontcolor}未安装${default_fontcolor}"
	else
		websocks_pid=$(ps -ef |grep "websocks" |grep -v "grep" | grep -v ".sh"| grep -v "init.d" |grep -v "service" |awk '{print $2}')
		if [[ ${websocks_pid} = "" ]]; then
			websocks_status="${red_fontcolor}未运行${default_fontcolor}"
		else
			websocks_status="${green_fontcolor}正在运行${default_fontcolor} | ${green_fontcolor}${websocks_pid}${default_fontcolor}"
		fi
	fi
	caddy_config=$(cat /usr/local/caddy/Caddyfile)
	if [[ ${caddy_config} = "" ]]; then
		caddy_status="${red_fontcolor}未安装${default_fontcolor}"
	else
		caddy_pid=$(ps -ef |grep "caddy" |grep -v "grep" | grep -v ".sh"| grep -v "init.d" |grep -v "service" |awk '{print $2}')
		if [[ ${caddy_pid} = "" ]]; then
			caddy_status="${red_fontcolor}未运行${default_fontcolor}"
		else
			caddy_status="${green_fontcolor}正在运行${default_fontcolor} | ${green_fontcolor}${caddy_pid}${default_fontcolor}"
		fi
	fi
}

function echo_install_list(){
	clear
	echo -e "脚本当前安装状态：${install_status}
--------------------------------------------------------------------------------------------------
安装Websocks:
	0.清除Caddy
	1.Websocks+TLS+网站伪装
--------------------------------------------------------------------------------------------------
Websocks当前运行状态：${websocks_status}
Caddy当前运行状态：${caddy_status}
	2.更新脚本
	3.更新程序
	4.卸载程序

	5.启动程序
	6.关闭程序
	7.重启程序
--------------------------------------------------------------------------------------------------
客户端运行指令：${websocks_start_command}
--------------------------------------------------------------------------------------------------"
	stty erase '^H' && read -p "请输入序号：" determine_type
	if [[ ${determine_type} = "" ]]; then
		clear
		echo -e "${error_font}请输入序号！"
		exit 1
	elif [[ ${determine_type} -lt 0 ]]; then
		clear
		echo -e "${error_font}请输入正确的序号！"
		exit 1
	elif [[ ${determine_type} -gt 7 ]]; then
		clear
		echo -e "${error_font}请输入正确的序号！"
		exit 1
	else
		data_processing
	fi
}

function data_processing(){
	clear
	echo -e "正在处理请求中..."
	if [[ ${determine_type} = "0" ]]; then
		uninstall_old
	elif [[ ${determine_type} = "2" ]]; then
		upgrade_shell_script
	elif [[ ${determine_type} = "3" ]]; then
		prevent_uninstall_check
		upgrade_program
		restart_service
	elif [[ ${determine_type} = "4" ]]; then
		prevent_uninstall_check
		uninstall_program
	elif [[ ${determine_type} = "5" ]]; then
		prevent_uninstall_check
		start_service
	elif [[ ${determine_type} = "6" ]]; then
		prevent_uninstall_check
		stop_service
	elif [[ ${determine_type} = "7" ]]; then
		prevent_uninstall_check
		restart_service
	else
		prevent_install_check
		os_update
		check_time
		generate_base_config
		clear
		echo -e "安装Websocks主程序中..."
		websocks_ver=$(wget -qO- "https://github.com/lzjluzijie/websocks/tags"|grep "/websocks/releases/tag/"|grep -v '\-apk'|head -n 1|awk -F "/tag/" '{print $2}'|sed 's/\">//')
		mkdir /usr/local/websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}建立文件夹成功。"
		else
			clear
			echo -e "${error_font}建立文件夹失败！"
			clear_install
			exit 1
		fi
		cd /usr/local/websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}进入文件夹成功。"
		else
			clear
			echo -e "${error_font}进入文件夹失败！"
			clear_install
			exit 1
		fi
		wget https://github.com/lzjluzijie/websocks/releases/download/${websocks_ver}/websocks_Linux_${system_bit}.tar.gz
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks下载成功。"
		else
			clear
			echo -e "${error_font}Websocks下载失败！"
			clear_install
			exit 1
		fi
		tar -xzf websocks_Linux_${system_bit}.tar.gz
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks解压成功。"
		else
			clear
			echo -e "${error_font}Websocks解压失败！"
			clear_install
			exit 1
		fi
		rm -rf LICENSE
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			clear_install
			exit 1
		fi
		rm -rf README.md
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			clear_install
			exit 1
		fi
		rm -rf README-zh.md
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			clear_install
			exit 1
		fi
		chmod 700 websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}设定Websocks权限成功。"
		else
			clear
			echo -e "${error_font}设定Websocks权限失败！"
			clear_install
			exit 1
		fi
		echo -e "${ok_font}Websocks安装成功。"
		if [[ ${determine_type} = "1" ]]; then
			clear
			echo -e "正在安装acme.sh中..."
			curl https://get.acme.sh | sh
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}acme.sh 安装成功。"
			else
				clear
				echo -e "${error_font}acme.sh 安装失败，请检查相关依赖是否正确安装。"
				clear_install
				exit 1
			fi
			bash <(curl https://raw.githubusercontent.com/ToyoDAdoubi/doubi/master/caddy_install.sh)
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}Caddy 安装成功。"
			else
				clear
				echo -e "${error_font}Caddy 安装失败，请检查相关依赖是否正确安装。"
				clear_install
				exit 1
			fi
			wget -O "/usr/local/caddy/Caddyfile" "https://raw.githubusercontent.com/1715173329/websocks-onekey/master/configs/websocks-tls-website.Caddyfile"
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}下载Caddy配置文件成功。"
			else
				clear
				echo -e "${error_font}下载Caddy配置文件失败！"
				clear_install
			fi
			clear
			install_port="443"
			check_port
			clear
			stty erase '^H' && read -p "请输入您的域名：" install_domain
			if [[ ${install_domain} = "" ]]; then
				clear
				echo -e "${error_font}请输入您的域名。"
				clear_install
				exit 1
			else
				clear
				echo -e "正在签发证书中..."
				bash ~/.acme.sh/acme.sh --issue -d ${install_domain} --standalone -k ec-256 --force
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}证书生成成功。"
					bash ~/.acme.sh/acme.sh --installcert -d ${install_domain} --fullchainpath /usr/local/websocks/pem.pem --keypath /usr/local/websocks/key.key --ecc
					if [[ $? -eq 0 ]];then
						clear
						echo -e "${ok_font}证书配置成功。"
					else
						clear
						echo -e "${error_font}证书配置失败！"
						clear_install
						exit 1
					fi
				else
					clear
					echo -e "${error_font}证书生成失败！"
					clear_install
					exit 1
				fi
				echo "${install_domain}" > /usr/local/websocks/full_domain.txt
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}Websocks 域名写入成功。"
				else
					clear
					echo -e "${error_font}Websocks 域名写入失败！"
					clear_install
					exit 1
				fi
				sed -i "s/PathUUID/${UUID}/g" "/usr/local/caddy/Caddyfile"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}Caddy UUID配置成功。"
				else
					clear
					echo -e "${error_font}Caddy UUID配置失败！"
					clear_install
					exit 1
				fi
				sed -i "s/WebsocksAddress/${install_domain}/g" "/usr/local/caddy/Caddyfile"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}Caddy 域名配置成功。"
				else
					clear
					echo -e "${error_font}Caddy 域名配置失败！"
					clear_install
					exit 1
				fi
				sed -i "s/WebsocksListenPort/${websocks_listen_port}/g" "/usr/local/caddy/Caddyfile"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}Caddy 监听端口配置成功。"
				else
					clear
					echo -e "${error_font}Caddy 监听端口配置失败！"
					clear_install
					exit 1
				fi
				cat <<-EOF > /etc/systemd/system/websocks.service
				[Unit]
				Description=websocks
				
				[Service]
				ExecStart=/usr/local/websocks/websocks server -l 127.0.0.1:${websocks_listen_port} -p /fuckgfw_gfwmotherfuckingboom/${UUID}
				Restart=always
				  
				[Install]
				WantedBy=multi-user.target
				EOF
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}写入Systemd配置成功。"
				else
					clear
					echo -e "${error_font}写入Systemd配置失败！"
					clear_install
					exit 1
				fi
				systemctl enable websocks.service
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}开启自启动成功。"
				else
					clear
					echo -e "${error_font}开启自启动失败！"
					clear_install
					exit 1
				fi
				mkdir /usr/local/websocks/pages
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}创建文件夹成功。"
				else
					clear
					echo -e "${error_font}创建文件夹失败！"
					clear_install
					exit 1
				fi
				cd /usr/local/websocks/pages
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}进入文件夹成功。"
				else
					clear
					echo -e "${error_font}进入文件夹失败！"
					clear_install
					exit 1
				fi
				wget -O "/usr/local/websocks/pages/websocks-page.zip" "https://github.com/1715173329/websocks-onekey/blob/master/pages/websocks-page.zip?raw=true"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}下载网页文件压缩包成功。"
				else
					clear
					echo -e "${error_font}下载网页文件压缩包失败！"
					clear_install
					exit 1
				fi
				unzip /usr/local/websocks/pages/websocks-page.zip
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}解压网页文件成功。"
				else
					clear
					echo -e "${error_font}解压网页文件失败！"
					clear_install
					exit 1
				fi
				rm -rf /usr/local/websocks/pages/websocks-page.zip
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}删除网页文件压缩包成功。"
				else
					clear
					echo -e "${error_font}删除网页文件压缩包失败！"
					clear_install
					exit 1
				fi
				sed -i "s/HTML_NUMBER/${html_number}/g" "/usr/local/websocks/pages/index.html"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}配置网页电话号成功。"
				else
					clear
					echo -e "${error_font}配置网页电话号失败！"
					clear_install
					exit 1
				fi
				sed -i "s/WebsocksAddress/${install_domain}/g" "/usr/local/websocks/pages/index.html"
				if [[ $? -eq 0 ]];then
					clear
					echo -e "${ok_font}配置网页域名成功。"
				else
					clear
					echo -e "${error_font}配置网页域名失败！"
					clear_install
					exit 1
				fi
			fi
			echo "1" > /usr/local/websocks/install_type.txt
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}写入安装信息成功。"
			else
				clear
				echo -e "${error_font}写入安装信息失败！"
				clear_install
				exit 1
			fi
			cd /root/
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}返回root文件夹成功。"
			else
				clear
				echo -e "${error_font}返回root文件夹失败！"
				clear_install
				exit 1
			fi
			restart_service
			echo_websocks_config
		fi
	fi
	echo -e "\n${ok_font}请求处理完毕。"
}

function uninstall_old(){
	clear
	echo -e "正在检查安装信息中..."
	clear
	stty erase '^H' && read -p "您是否需要卸载Caddy？[Y/N,Default:N]" uninstall_caddy_right
	if [[ ${uninstall_caddy_right} == [Yy] ]]; then
		if [[ ${caddy_status} = "${red_fontcolor}未安装${default_fontcolor}" ]]; then
			clear
			echo -e "${error_font}您未安装Caddy。"
		else
			service caddy stop
			update-rc.d -f caddy remove
			rm -rf /etc/init.d/caddy
			rm -rf /root/.caddy
			rm -rf /usr/local/caddy
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}Caddy卸载成功。"
			else
				clear
				echo -e "${error_font}Caddy卸载失败！"
			fi
		fi
	else
		clear
		echo -e "${ok_font}取消卸载Caddy成功。"
	fi
}

function upgrade_shell_script(){
	clear
	echo -e "正在更新脚本中..."
	filepath=$(cd "$(dirname "$0")"; pwd)
	filename=$(echo -e "${filepath}"|awk -F "$0" '{print $1}')
	curl -O https://raw.githubusercontent.com/1715173329/websocks-onekey/master/websocks-go.sh
	if [[ $? -eq 0 ]];then
		clear
		echo -e "${ok_font}脚本更新成功，脚本位置：\"${green_backgroundcolor}${filename}/websocks-go.sh${default_fontcolor}\"，使用：\"${green_backgroundcolor}bash ${filename}/websocks-go.sh${default_fontcolor}\"。"
	else
		clear
		echo -e "${error_font}脚本更新失败！"
	fi
}

function prevent_uninstall_check(){
	clear
	echo -e "正在检查安装状态中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [ "${install_type}" = "" ]; then
		clear
		echo -e "${error_font}您未安装本程序。"
		exit 1
	else
		echo -e "${ok_font}您已安装本程序，正在执行相关命令中..."
	fi
}

function start_service(){
	clear
	echo -e "正在启动服务中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [ "${install_type}" -eq "1" ]; then
		if [[ ${websocks_pid} -eq 0 ]]; then
			service websocks start
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}Websocks 启动成功。"
			else
				clear
				echo -e "${error_font}Websocks 启动失败！"
			fi
		else
			clear
			echo -e "${error_font}Websocks 正在运行。"
		fi
		if [[ ${caddy_pid} -eq 0 ]]; then
			service caddy start
			if [[ $? -eq 0 ]];then
				echo -e "${ok_font}Caddy 启动成功。"
				exit 0
			else
				echo -e "${error_font}Caddy 启动失败！"
				exit 1
			fi
		else
			echo -e "${error_font}Caddy 正在运行。"
			exit 1
		fi
	fi
}

function stop_service(){
	clear
	echo -e "正在停止服务中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [ "${install_type}" -eq "1" ]; then
		if [[ ${websocks_pid} -eq 0 ]]; then
			clear
			echo -e "${error_font}Websocks 未在运行。"
		else
			service websocks stop
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}Websocks 停止成功。"
			else
				clear
				echo -e "${error_font}Websocks 停止失败！"
			fi
		fi
		if [[ ${caddy_pid} -eq 0 ]]; then
			echo -e "${error_font}Caddy 未在运行。"
		else
			service caddy stop
			if [[ $? -eq 0 ]];then
				echo -e "${ok_font}Caddy 停止成功。"
				exit 0
			else
				echo -e "${error_font}Caddy 停止失败！"
				exit 1
			fi
		fi
	fi
}

function restart_service(){
	clear
	echo -e "正在重启服务中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [ "${install_type}" -eq "1" ]; then
		service websocks restart
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks 重启成功。"
		else
			clear
			echo -e "${error_font}Websocks 重启失败！"
		fi
		service caddy restart
		if [[ $? -eq 0 ]];then
			echo -e "${ok_font}Caddy 重启成功。"
		else
			echo -e "${error_font}Caddy 重启失败！"
		fi
	fi
}

function prevent_install_check(){
	clear
	echo -e "正在检测安装状态中..."
	if [[ ${determine_type} -lt 9 ]]; then
		if [[ ${install_status} = "${green_fontcolor}已安装${default_fontcolor}" ]]; then
			echo -e "${error_font}您已经安装过了，请勿再次安装，若您需要切换至其他模式，请先卸载后再使用安装功能。"
			exit 1
		elif [[ ${websocks_status} = "${red_fontcolor}未安装${default_fontcolor}" ]]; then
			if [[ ${determine_type} -lt 8 ]]; then
				echo -e "${ok_font}检测完毕，符合要求，正在执行命令中..."
			else
				if [[ ${caddy_status} = "${red_fontcolor}未安装${default_fontcolor}" ]]; then
					echo -e "${ok_font}检测完毕，符合要求，正在执行命令中..."
				else
					echo -e "${error_font}您的VPS上已经安装Caddy，请勿再次安装，若您需要使用本脚本，请先卸载后再使用安装功能。"
					exit 1
				fi
			fi
		else
			echo -e "${error_font}您的VPS上已经安装Websocks，请勿再次安装，若您需要使用本脚本，请先卸载后再使用安装功能。"
			exit 1
		fi
	fi
}

function uninstall_program(){
	clear
	echo -e "正在卸载中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [[ "${install_type}" -eq "1" ]]; then
		full_domain=$(cat /usr/local/websocks/full_domain.txt)
		bash ~/.acme.sh/acme.sh --revoke -d ${full_domain} --ecc
		bash ~/.acme.sh/acme.sh --remove -d ${full_domain} --ecc
		rm -rf ~/.acme.sh
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Acme卸载成功。"
		else
			clear
			echo -e "${error_font}Acme卸载失败！"
		fi
		service websocks stop
		systemctl disable websocks.service
		rm -rf /etc/systemd/system/websocks.service
		update-rc.d -f websocks remove
		rm -rf /usr/local/websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks卸载成功。"
		else
			clear
			echo -e "${error_font}Websocks卸载失败！"
		fi
		service caddy stop
		update-rc.d -f caddy remove
		rm -rf /etc/init.d/caddy
		rm -rf /root/.caddy
		rm -rf /usr/local/caddy
		if [[ $? -eq 0 ]];then
			echo -e "${ok_font}Caddy卸载成功。"
		else
			echo -e "${error_font}Caddy卸载失败！"
		fi
	fi
}

function upgrade_program(){
	clear
	echo -e "正在更新程序中..."
	install_type=$(cat /usr/local/websocks/install_type.txt)
	if [ "${install_type}" -eq "1" ]; then
		cd /usr/local/websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}进入文件夹成功。"
		else
			clear
			echo -e "${error_font}进入文件夹失败！"
			exit 1
		fi
		mv websocks websocks.bak
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}备份旧文件成功。"
		else
			clear
			echo -e "${error_font}备份旧文件失败！"
			exit 1
		fi
		wget https://github.com/lzjluzijie/websocks/releases/download/${websocks_ver}/websocks_Linux_${system_bit}.tar.gz
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks下载成功。"
			rm -rf websocks.bak
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}删除备份文件成功。"
			else
				clear
				echo -e "${error_font}删除备份文件失败！"
				exit 1
			fi
		else
			clear
			echo -e "${error_font}Websocks下载失败！"
			mv websocks.bak websocks
			if [[ $? -eq 0 ]];then
				clear
				echo -e "${ok_font}恢复备份文件成功。"
			else
				clear
				echo -e "${error_font}恢复备份文件失败！"
				exit 1
			fi
			exit 1
		fi
		tar -xzf websocks_Linux_${system_bit}.tar.gz
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks解压成功。"
		else
			clear
			echo -e "${error_font}Websocks解压失败！"
			exit 1
		fi
		rm -rf LICENSE
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			exit 1
		fi
		rm -rf README.md
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			exit 1
		fi
		rm -rf README-zh.md
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}删除无用文件成功。"
		else
			clear
			echo -e "${error_font}删除无用文件失败！"
			exit 1
		fi
		chmod 700 websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}设定Websocks权限成功。"
		else
			clear
			echo -e "${error_font}设定Websocks权限失败！"
			exit 1
		fi
		echo -e "${ok_font}Websocks 更新成功。"
		bash <(curl https://raw.githubusercontent.com/ToyoDAdoubi/doubi/master/caddy_install.sh)
		if [[ $? -eq 0 ]];then
			echo -e "${ok_font}Caddy 更新成功。"
		else
			echo -e "${error_font}Caddy 更新失败！"
		fi
	fi
}

function clear_install(){
	clear
	echo -e "正在卸载中..."
	if [ "${determine_type}" -eq "1" ]; then
		full_domain=$(cat /usr/local/websocks/full_domain.txt)
		bash ~/.acme.sh/acme.sh --revoke -d ${full_domain} --ecc
		bash ~/.acme.sh/acme.sh --remove -d ${full_domain} --ecc
		rm -rf ~/.acme.sh
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Acme卸载成功。"
		else
			clear
			echo -e "${error_font}Acme卸载失败！"
		fi
		service websocks stop
		systemctl disable websocks.service
		rm -rf /etc/systemd/system/websocks.service
		update-rc.d -f websocks remove
		rm -rf /usr/local/websocks
		if [[ $? -eq 0 ]];then
			clear
			echo -e "${ok_font}Websocks卸载成功。"
		else
			clear
			echo -e "${error_font}Websocks卸载失败！"
		fi
		service caddy stop
		update-rc.d -f caddy remove
		rm -rf /etc/init.d/caddy
		rm -rf /root/.caddy
		rm -rf /usr/local/caddy
		if [[ $? -eq 0 ]];then
			echo -e "${ok_font}Caddy卸载成功。"
		else
			echo -e "${error_font}Caddy卸载失败！"
		fi
	fi
}

function os_update(){
	clear
	echo -e "正在安装/更新系统组件中..."
	clear
	apt-get -y update
	apt-get -y upgrade
	apt-get -y install wget curl ntpdate unzip socat lsof cron iptables
	if [[ $? -ne 0 ]];then
		clear
		echo -e "${error_font}系统组件更新失败！"
		exit 1
	else
		clear
		echo -e "${ok_font}系统组件更新成功。"
	fi
}

function check_time(){
	clear
	echo -e "正在对时中..."
	rm -rf /etc/localtime
	cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
	ntpdate time.nist.gov
	if [[ $? -eq 0 ]];then
		clear
		echo -e "${ok_font}时间同步成功。"
		echo -e "${ok_font}当前系统时间 $(date -R) （请注意时区间时间换算，换算后时间误差应为三分钟以内）"
	else
		clear
		echo -e "${error_font}时间同步失败，请检查ntpdate服务是否正常工作。"
		echo -e "${error_font}当前系统时间 $(date -R) ，如果和你的本地时间有误差，请手动调整。"
	fi 
}

function generate_base_config(){
	clear
	echo "正在生成基础信息中..."
	hostname=$(hostname)
	Address=$(curl https://ipinfo.io/ip)
	UUID=$(cat /proc/sys/kernel/random/uuid)
	let websocks_listen_port=$RANDOM+10000
	let html_number=$RANDOM+10000
	if [[ ${Address} = "" ]]; then
		clear
		echo -e "${error_font}读取vps_ip失败！"
		exit 1
	elif [[ ${UUID} = "" ]]; then
		clear
		echo -e "${error_font}生成UUID失败！"
		exit 1
	elif [[ ${websocks_listen_port} = "" ]]; then
		clear
		echo -e "${error_font}生成websocks监听端口失败！"
		exit 1
	elif [[ ${html_number} = "" ]]; then
		clear
		echo -e "${error_font}生成网页随机数字失败！"
		exit 1
	else
		clear
		echo -e "${ok_font}您的vps_ip为：${Address}"
		echo -e "${ok_font}生成的UUID为：${UUID}"
		echo -e "${ok_font}生成的Websocks监听端口为：${websocks_listen_port}"
		echo -e "${ok_font}生成的网页随机数字为：${html_number}"
	fi
}

function check_port(){
	clear
	echo "正在检查端口占用情况："
	if [[ 0 -eq $(lsof -i:"${install_port}" | wc -l) ]];then
		clear
		echo -e "${ok_font}端口未被占用。"
		open_port
	else
		clear
		echo -e "${error_font}端口被占用，请切换使用其他端口。"
		clear_install
		exit 1
	fi
}

function open_port(){
	clear
	echo -e "正在设置防火墙中..."
	iptables-save > /etc/iptables.up.rules
	echo -e '#!/bin/bash\n/sbin/iptables-restore < /etc/iptables.up.rules' > /etc/network/if-pre-up.d/iptables
	chmod +x /etc/network/if-pre-up.d/iptables
	iptables -I INPUT -m state --state NEW -m tcp -p tcp --dport ${install_port} -j ACCEPT
	iptables -I INPUT -m state --state NEW -m udp -p udp --dport ${install_port} -j ACCEPT
	iptables-save > /etc/iptables.up.rules
	if [[ $? -eq 0 ]];then
		clear
		echo -e "${ok_font}端口开放配置成功。"
	else
		clear
		echo -e "${error_font}端口开放配置失败！"
		clear_install
		exit 1
	fi
}

function echo_websocks_config(){
	if [[ ${determine_type} = "1" ]]; then
		clear
		run_command="./websocks client -l 127.0.0.1:1080 -s wss://${install_domain}/fuckgfw_gfwmotherfuckingboom/${UUID}" 
		echo -e "您的连接信息如下："
		echo -e "WSS地址：${install_domain}"
		echo -e "端口(Port)：${install_port}"
		echo -e "WSS目录：/fuckgfw_gfwmotherfuckingboom/${UUID}"
		echo -e "客户端运行指令：${green_backgroundcolor}${run_command}${default_fontcolor}"
	fi
	echo -e "${run_command}" > /usr/local/websocks/run_command.txt
}

function main(){
	set_fonts_colors
	check_os
	check_install_status
	echo_install_list
}

	main
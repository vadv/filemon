%define version unknown
%define bin_name gofilemon
%define debug_package %{nil}

Name:           %{bin_name}
Version:        %{version}
Release:        1%{?dist}
Summary:        log file parser
License:        BSD
URL:            http://git.itv.restr.im/infra/%{bin_name}
Source:         %{bin_name}-%{version}.tar.gz
BuildRequires:  make
BuildRequires:  go

%define restream_dir /opt/restream/
%define restream_zabbix_bin_dir %{restream_dir}/zabbix/bin

%description
This package provides log parser.

%prep
%setup

%build
make

%install
mkdir -p %{buildroot}%{restream_zabbix_bin_dir}
install bin/%{bin_name} %{buildroot}%{restream_zabbix_bin_dir}

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_zabbix_bin_dir}/%{bin_name}

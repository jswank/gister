Name:		gister
Version:	1.0.0
Release:	1%{?dist}
Summary:	Command line gist generator
License:	MIT
URL:		https://github.com/jswank/gister
Source0:	%{name}-%{version}.tar.bz2
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:	asciidoc, xmlto

%description
The gister program creates Github gists using one or more files and
outputs the URL of the gist created.  Github OAuth authentication must
be used.

%prep
%setup -q

%build
make %{?_smp_mflags}

%install
rm -rf %{buildroot}
make install DESTDIR=%{buildroot}

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{_bindir}/*
%{_mandir}/man1/*

%changelog
* Sun Jan 13 2013 Jason Swank <jswank@scalene.net> 1.0.0-1
- intial packaging
